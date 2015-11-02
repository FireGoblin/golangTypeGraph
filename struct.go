package main

import (
	"github.com/firegoblin/gographviz"
	"go/ast"
)

//A node is responsible for the incoming edges to it

//A node type
//implements gographviz.GraphableNode
type structNode struct {
	target *baseType

	//if this is not nil, it is a redefined type
	//edge is drawable if parent.base.node is not nil
	parent *Type

	//fields should only be empty when a redefined type
	fields []namedType

	//receiver functions on this type
	receiverFunctions []receiverFunction

	//structs or interfaces included anonymously in this struct
	inheritedTypes []*baseType

	//interfaces this node implements
	interfaceCache []*interfaceNode

	//any attrs need for drawing in the graph
	extraAttrs gographviz.Attrs

	//either StructType or parent type for embedded type
	astNode ast.Expr
}

func (s *structNode) addFunction(f *function, field *ast.Field) {
	s.receiverFunctions = append(s.receiverFunctions, newReceiverFunction(f, field))
	f.isReceiver = true
}

func (s *structNode) String() string {
	return s.target.String()
}

func (s *structNode) Name() string {
	return s.target.Name()
}

func (s *structNode) Attrs() gographviz.Attrs {
	retval := make(map[string]string)
	retval["shape"] = "record"
	retval["label"] = s.label()
	return retval
}

func containsNode(list []gographviz.GraphableNode, g gographviz.GraphableNode) bool {
	for _, v := range list {
		if g == v {
			return true
		}
	}
	return false
}

func (s *structNode) Edges() []*gographviz.Edge {
	parentEdge := s.parentEdge()
	var retval []*gographviz.Edge

	if parentEdge != nil {
		retval = make([]*gographviz.Edge, 0, len(s.inheritedTypes)+1)
	} else {
		retval = make([]*gographviz.Edge, 0, len(s.inheritedTypes))
	}

	for _, v := range s.inheritedTypes {
		//TODO: decide on attrs
		if _, ok := v.node.(*unknownNode); !ok && v.node != nil {
			retval = append(retval, &gographviz.Edge{v.node.Name(), "", s.Name(), "", true, inheritedAttrs()})
		}
	}

	fieldList := make([]gographviz.GraphableNode, len(s.fields))
	for _, f := range s.fields {
		holder := f.node()
		if _, ok := holder.(*unknownNode); ok {
			continue
		}
		if holder != nil {
			//avoid duplicates
			if !containsNode(fieldList, holder) {
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

func (s *structNode) label() string {
	retval := "\"{" + s.String() + "|"

	if s.parent != nil {
		retval += s.parent.stringRelativePkg(s.target.pkgName)
	}

	for _, v := range s.inheritedTypes {
		retval += v.stringRelativePkg(s.target.pkgName) + "\\l"
	}

	for _, v := range s.fields {
		retval += v.stringRelativePkg(s.target.pkgName) + "\\l"
	}

	retval += "|"

	for _, v := range s.receiverFunctions {
		retval += v.SlimString() + "\\l"
	}

	retval += "}\""

	return retval
}

func (s *structNode) parentEdge() *gographviz.Edge {
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

//no mutation
func (s *structNode) allreceiverFunctions() []*function {
	retval := make([]*function, len(s.receiverFunctions))
	for _, v := range s.receiverFunctions {
		retval = append(retval, v.f)
	}

	for _, v := range s.inheritedTypes {
		switch w := v.node.(type) {
		case *structNode:
			retval = append(retval, w.allreceiverFunctions()...)
		case *interfaceNode:
			retval = append(retval, w.allRequiredFunctions()...)
		}
	}

	return retval
}

//no mutation
func (s *structNode) implementsInterface(i *interfaceNode) bool {
	required := i.allRequiredFunctions()
	have := s.allreceiverFunctions()

	return containsAll(have, required)
}

func (s *structNode) setInterfacesImplemented(i []*interfaceNode) {
	retval := s.interfacesImplemented(i)
	s.interfaceCache = make([]*interfaceNode, len(retval))
	copy(s.interfaceCache, retval)
}

//filter pattern
func (s *structNode) interfacesImplemented(i []*interfaceNode) []*interfaceNode {
	//cache call
	if i == nil {
		return s.interfaceCache
	}

	var retval []*interfaceNode

	for _, v := range i {
		if s.implementsInterface(v) {
			retval = append(retval, v)
		}
	}

	return retval
}

func (s *structNode) remakeStructInternals(spec *ast.TypeSpec) {
	switch t := spec.Type.(type) {
	case *ast.StructType:
		//struct
		FlattenedFields := Flattened(t.Fields)
		for _, v := range FlattenedFields.List {
			if len(v.Names) != 0 {
				s.fields = append(s.fields, newNamedTypeFromField(v))
			} else {
				lookup := typeMap.lookupOrAddFromExpr(v.Type)
				if lookup.base.node != nil {
					s.inheritedTypes = append(s.inheritedTypes, lookup.base)
				} else {
					s.inheritedTypes = append(s.inheritedTypes, newUnknown(s, lookup.base).target)
				}
			}
		}
	default:
		//redefined type
		s.parent = typeMap.lookupOrAddFromExpr(t)
	}
}

//possibilities for lines:
//Type -> inheritedStruct
//(comma seperated list of names) Type -> namedTypes
//b: the baseType for this struct
//lines: lines from the structs declaration block, preceeding and trailing whitespace removed
func newStruct(spec *ast.TypeSpec, b *baseType) *structNode {
	//should only be used with declarations, if struct is in field names use newStructUnknown
	retval := &structNode{b, nil, make([]namedType, 0), make([]receiverFunction, 0), make([]*baseType, 0), nil, nil, spec.Type}

	retval.remakeStructInternals(spec)

	b.addNode(retval)
	return retval
}
