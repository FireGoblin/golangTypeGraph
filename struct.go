package main

//import . "regexp"
//import "strings"
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
	receiverFunctions []*Function

	//structs or interfaces included anonymously in this struct
	inheritedStructs []gographviz.GraphableNode

	//structs this type is included in anonymously
	includedIn []*Struct

	//interfaces this node implements
	interfaceCache []*Interface

	//any attrs need for drawing in the graph
	extraAttrs gographviz.Attrs

	//either StructType or Ident for embedded type
	astNode ast.Expr
}

func (s *Struct) String() string {
	return s.target.name
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
		retval = make([]*gographviz.Edge, 0, len(s.inheritedStructs)+1)
	} else {
		retval = make([]*gographviz.Edge, 0, len(s.inheritedStructs))
	}

	for _, v := range s.inheritedStructs {
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

	for _, v := range s.inheritedStructs {
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
	s.interfaceCache = make([]*Interface, 0, len(retval))
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
	return len(s.inheritedStructs) > 0
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

func makeRedefinedTyp(b *BaseType, s string) *Struct {
	retval := Struct{b, typeMap.lookupOrAdd(s), make([]NamedType, 0), make([]*Function, 0), make([]gographviz.GraphableNode, 0), make([]*Struct, 0), nil, nil, nil}
	b.addNode(&retval)

	return &retval
}

//for if struct is found as an Anonymous member of something else first
func makeStructUnknown(b *BaseType, source *Struct) *Struct {
	retval := Struct{b, nil, make([]NamedType, 0), make([]*Function, 0), make([]gographviz.GraphableNode, 0), make([]*Struct, 0), nil, nil, nil}
	b.addNode(&retval)

	retval.includedIn = append(retval.includedIn, source)

	return &retval
}

//possibilities for lines:
//Type -> inheritedStruct
//(comma seperated list of names) Type -> NamedTypes
//b: the baseType for this struct
//lines: lines from the structs declaration block, preceeding and trailing whitespace removed
func makeStruct(spec *ast.TypeSpec, typ *BaseType) *Struct {
	//should only be used with declarations, if struct is in field names use makeStructUnknown
	// b := makeType(spec.Name)
	// retval := Struct{b, nil, make([]NamedType, 0), make([]*Function, 0), make([]*Struct, 0), make([]*Struct, 0)}
	// b.addNode(&retval)

	// switch typ := spec.Type.(type) {
	// case *ast.StructType:
	// 	//struct
	// case *ast.Ident:
	// 	//redefined type
	// 	retval.parent = typeMap.lookupOrAdd(typ.Name)
	// default:
	// 	panic("unexpected type in makeStruct")
	// }

	// for _, v := range lines {
	// 	ntp := NamedTypeParser.FindStringSubmatch(v)
	// 	if len(ntp) != 0 {
	// 		str := strings.Split(ntp[1], ", ")
	// 		typ := typeMap.lookupOrAdd(ntp[2])
	// 		for _, s := range str {
	// 			retval.fields = append(retval.fields, NamedType{s, typ})
	// 		}
	// 	} else {
	// 		typ := typeMap.lookupOrAdd(v)
	// 		var struc *Struct

	// 		if typ.base.node == nil {
	// 			struc := makeStructUnknown(typ.base, &retval)
	// 			typ.base.node = struc
	// 		} else {
	// 			var ok bool
	// 			struc, ok = typ.base.node.(*Struct)
	// 			if !ok {
	// 				panic("Could not find struct of anonymous member")
	// 			}
	// 		}
	// 		retval.inheritedStructs = append(retval.inheritedStructs, struc)
	// 	}
	// }

	// return &retval
	return nil
}
