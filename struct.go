package main

//import . "regexp"
//import "strings"
import "go/ast"
import "github.com/firegoblin/gographviz"

//import "fmt"

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
	receiverFunctions []*Function

	//structs or interfaces included anonymously in this struct
	inheritedTypes []gographviz.GraphableNode

	//structs this type is included in anonymously
	includedIn []*Struct

	//interfaces this node implements
	interfaceCache []*Interface

	//any attrs need for drawing in the graph
	extraAttrs gographviz.Attrs

	//either StructType or Ident for embedded type
	astNode ast.Expr
}

func (s *Struct) AddFunction(f *Function) {
	s.receiverFunctions = append(s.receiverFunctions, f)
}

func (s *Struct) String() string {
	retval := "Struct: " + s.target.name + "\n"
	retval += "Fields:\n"
	for _, v := range s.fields {
		retval += v.String() + "\n"
	}
	retval += "Receiver Functions:\n"
	for _, v := range s.receiverFunctions {
		retval += v.String() + "\n"
	}

	return retval
}

func (s *Struct) Name() string {
	return s.target.name
}

//TODO: improve
func (s *Struct) Attrs() gographviz.Attrs {
	return nil
}

func (s *Struct) parentEdge() *gographviz.Edge {
	if s.parent == nil {
		return nil
	}

	parentNode := s.parent.base.node
	//TODO: better attrs
	return &gographviz.Edge{parentNode.Name(), "", s.Name(), "", true, nil}
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
		retval = append(retval, &gographviz.Edge{v.Name(), "", s.Name(), "", true, nil})
	}

	if parentEdge != nil {
		retval = append(retval, parentEdge)
	}
	return retval
}

//no mutation
func (s *Struct) allReceiverFunctions() []*Function {
	retval := make([]*Function, len(s.receiverFunctions))
	c := copy(retval, s.receiverFunctions)
	if c != len(s.receiverFunctions) {
		panic("copy failed in allRequiredFunctions")
	}

	for _, v := range s.inheritedTypes {
		switch w := v.(type) {
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

func (s *Struct) addToIncludedIn(x *Struct) {
	s.includedIn = append(s.includedIn, x)
}

//for if struct is found as an Anonymous member of something else first
func makeStructUnknown(source *Struct, b *BaseType) *Struct {
	retval := &Struct{b, nil, make([]NamedType, 0), make([]*Function, 0), make([]gographviz.GraphableNode, 0), make([]*Struct, 0), nil, nil, nil}
	b.addNode(retval)

	retval.includedIn = append(retval.includedIn, source)

	return retval
}

func (s *Struct) remakeStruct(spec *ast.TypeSpec) *Struct {
	switch t := spec.Type.(type) {
	case *ast.StructType:
		//struct
		flattenedFields := flattened(t.Fields)
		for _, v := range flattenedFields.List {
			if len(v.Names) != 0 {
				s.fields = append(s.fields, NamedType{v.Names[0].Name, typeMap.lookupOrAdd(String(v.Type))})
			} else {
				lookup := typeMap.lookupOrAdd(String(v.Type))
				node := lookup.base.node
				if node != nil {
					s.inheritedTypes = append(s.inheritedTypes, node)
					node.(*Struct).addToIncludedIn(s)
				} else {
					s.inheritedTypes = append(s.inheritedTypes, makeStructUnknown(s, lookup.base))
				}
			}
		}
	case *ast.Ident:
		//redefined type
		s.parent = typeMap.lookupOrAdd(t.Name)
	default:
		panic("unexpected type in makeStruct")
	}

	return s
}

//possibilities for lines:
//Type -> inheritedStruct
//(comma seperated list of names) Type -> NamedTypes
//b: the baseType for this struct
//lines: lines from the structs declaration block, preceeding and trailing whitespace removed
func makeStruct(spec *ast.TypeSpec, b *BaseType) *Struct {
	//should only be used with declarations, if struct is in field names use makeStructUnknown
	retval := &Struct{b, nil, make([]NamedType, 0), make([]*Function, 0), make([]gographviz.GraphableNode, 0), make([]*Struct, 0), nil, nil, spec.Type}

	switch t := spec.Type.(type) {
	case *ast.StructType:
		//struct
		flattenedFields := flattened(t.Fields)
		for _, v := range flattenedFields.List {
			if len(v.Names) != 0 {
				retval.fields = append(retval.fields, NamedType{v.Names[0].Name, typeMap.lookupOrAdd(String(v.Type))})
			} else {
				lookup := typeMap.lookupOrAdd(String(v.Type))
				node := lookup.base.node
				if node != nil {
					retval.inheritedTypes = append(retval.inheritedTypes, node)
					node.(*Struct).addToIncludedIn(retval)
				} else {
					retval.inheritedTypes = append(retval.inheritedTypes, makeStructUnknown(retval, lookup.base))
				}
			}
		}
	case *ast.Ident:
		//redefined type
		retval.parent = typeMap.lookupOrAdd(t.Name)
	default:
		panic("unexpected type in makeStruct")
	}

	//fmt.Println("makeStruct:")
	//fmt.Println(retval)
	b.addNode(retval)
	return retval
}
