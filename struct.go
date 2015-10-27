package main

import "go/ast"
import "github.com/firegoblin/gographviz"

//A node is responsible for the incoming edges to it

//A node type
//implements gographviz.GraphableNode
type Struct struct {
	target *BaseType

	//if this is not nil, it is a redefined type
	//edge is drawable if parent.base.node is not nil
	parent *Type

	//fields should only be empty when a redefined type
	fields []NamedType

	//receiver functions on this type
	receiverFunctions []ReceiverFunction

	//structs or interfaces included anonymously in this struct
	inheritedTypes []*BaseType

	//interfaces this node implements
	interfaceCache []*Interface

	//any attrs need for drawing in the graph
	extraAttrs gographviz.Attrs

	//either StructType or parent type for embedded type
	astNode ast.Expr
}

func (s *Struct) AddFunction(f *Function, field *ast.Field) {
	s.receiverFunctions = append(s.receiverFunctions, NewReceiverFunction(f, field))
	f.isReceiver = true
}

func (s *Struct) String() string {
	return s.target.String()
}

func (s *Struct) Name() string {
	return s.target.Name()
}

func (s *Struct) label() string {
	retval := "\"{" + s.String() + "|"

	if s.parent != nil {
		retval += s.parent.StringRelativePkg(s.target.pkgName)
	}

	for _, v := range s.inheritedTypes {
		retval += v.String() + "\\l"
	}

	for _, v := range s.fields {
		retval += v.StringRelativePkg(s.target.pkgName) + "\\l"
	}

	retval += "|"

	for _, v := range s.receiverFunctions {
		retval += v.SlimString() + "\\l"
	}

	retval += "}\""

	return retval
}

//TODO: improve
func (s *Struct) Attrs() gographviz.Attrs {
	retval := make(map[string]string)
	retval["shape"] = "record"
	retval["label"] = s.label()
	return retval
}

func (s *Struct) parentEdge() *gographviz.Edge {
	if s.parent == nil || s.parent.base.node == nil {
		return nil
	}

	//TODO: better handling of derivative types
	//TODO: better attrs
	return &gographviz.Edge{s.parent.base.node.Name(), "", s.Name(), "", true, parentAttrs()}
}

func inheritedAttrs() map[string]string {
	retval := make(map[string]string)
	retval["label"] = "inherited"
	retval["style"] = "solid"
	return retval
}

func parentAttrs() map[string]string {
	retval := make(map[string]string)
	retval["label"] = "parent"
	retval["style"] = "solid"
	return retval
}

func fieldAttrs() map[string]string {
	retval := make(map[string]string)
	retval["label"] = "field"
	retval["style"] = "dashed"
	return retval
}

//TODO: add parent edge
func (s *Struct) Edges() []*gographviz.Edge {
	parentEdge := s.parentEdge()
	var retval []*gographviz.Edge

	if parentEdge != nil {
		retval = make([]*gographviz.Edge, 0, len(s.inheritedTypes)+1)
	} else {
		retval = make([]*gographviz.Edge, 0, len(s.inheritedTypes))
	}

	for _, v := range s.inheritedTypes {
		//TODO: decide on attrs
		retval = append(retval, &gographviz.Edge{v.node.Name(), "", s.Name(), "", true, inheritedAttrs()})
	}

	fieldList := make([]gographviz.GraphableNode, len(s.fields))
	for _, f := range s.fields {
		holder := f.Node()
		if holder != nil {
			//avoid duplicates
			found := false
			for _, v := range fieldList {
				if holder == v {
					found = true
					break
				}
			}
			if !found {
				retval = append(retval, &gographviz.Edge{holder.Name(), "", s.Name(), "", true, fieldAttrs()})
				fieldList = append(fieldList, holder)
			}
		}
	}

	if parentEdge != nil {
		retval = append(retval, parentEdge)
	}
	return retval
}

//no mutation
func (s *Struct) allReceiverFunctions() []*Function {
	retval := make([]*Function, len(s.receiverFunctions))
	for _, v := range s.receiverFunctions {
		retval = append(retval, v.f)
	}

	for _, v := range s.inheritedTypes {
		switch w := v.node.(type) {
		case *Struct:
			retval = append(retval, w.allReceiverFunctions()...)
		case *Interface:
			retval = append(retval, w.allRequiredFunctions()...)
		}
	}

	return retval
}

//no mutation
func (s *Struct) implementsInterface(i *Interface) bool {
	required := i.allRequiredFunctions()
	have := s.allReceiverFunctions()

	for _, v := range required {
		found := false
		for _, j := range have {
			if j == v {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func (s *Struct) setInterfacesImplemented(i []*Interface) {
	retval := s.interfacesImplemented(i)
	s.interfaceCache = make([]*Interface, len(retval))
	copy(s.interfaceCache, retval)
}

//filter pattern
func (s *Struct) interfacesImplemented(i []*Interface) []*Interface {
	//cache call
	if i == nil {
		return s.interfaceCache
	}

	var retval []*Interface

	for _, v := range i {
		if s.implementsInterface(v) {
			retval = append(retval, v)
		}
	}

	return retval
}

func (s *Struct) isComposite() bool {
	return len(s.inheritedTypes) > 0
}

func (s *Struct) isRedefinedType() bool {
	return s.parent != nil
}

func (s *Struct) getFields() []NamedType {
	if s.parent == nil {
		return s.fields
	}

	n := s.parent.base.node
	if n == nil {
		return nil
	}

	struc, ok := n.(*Struct)
	if ok {
		return struc.fields
	}

	return nil
}

//for if struct is found as an Anonymous member of something else first
func makeUnknown(source *Struct, b *BaseType) *Unknown {
	retval := &Unknown{b}
	b.addNode(retval)

	return retval
}

//possibilities for lines:
//Type -> inheritedStruct
//(comma seperated list of names) Type -> NamedTypes
//b: the baseType for this struct
//lines: lines from the structs declaration block, preceeding and trailing whitespace removed
func makeStruct(spec *ast.TypeSpec, b *BaseType) *Struct {
	//should only be used with declarations, if struct is in field names use makeStructUnknown
	retval := &Struct{b, nil, make([]NamedType, 0), make([]ReceiverFunction, 0), make([]*BaseType, 0), nil, nil, spec.Type}

	switch t := spec.Type.(type) {
	case *ast.StructType:
		//struct
		flattenedFields := flattened(t.Fields)
		for _, v := range flattenedFields.List {
			if len(v.Names) != 0 {
				retval.fields = append(retval.fields, NamedTypeFromField(v))
			} else {
				lookup := typeMap.lookupOrAddFromExpr(v.Type)
				if lookup.base.node != nil {
					retval.inheritedTypes = append(retval.inheritedTypes, lookup.base)
				} else {
					retval.inheritedTypes = append(retval.inheritedTypes, makeUnknown(retval, lookup.base).target)
				}
			}
		}
	default:
		//redefined type
		retval.parent = typeMap.lookupOrAddFromExpr(t)
	}

	b.addNode(retval)
	return retval
}
