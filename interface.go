package main

import "go/ast"
import "github.com/firegoblin/gographviz"

import "fmt"

//import . "regexp"

//A node type
//implements gographviz.GraphableNode
type Interface struct {
	target *BaseType

	//any functions required to implement the interface
	//does not include functions inherited indirectly through inheritedInterfaces
	requiredFunctions []*Function

	//interfaces this inherits from
	//if not zero, is composite interface
	//Draw edges from these
	inheritedInterfaces []*Interface

	//interfaces this is included in
	//TODO: change to allow for being included in structs
	includedIn []gographviz.GraphableNode

	//structs
	//Draw edges from these
	implementedByCache []*Struct

	extraAttrs gographviz.Attrs

	astNode *ast.InterfaceType
}

func (i *Interface) String() string {
	retval := "Interface: " + i.target.name + "\n"

	retval += "Receiver Functions:\n"
	for _, v := range i.requiredFunctions {
		retval += v.String() + "\n"
	}

	return retval
}

func (i *Interface) Name() string {
	return i.target.name
}

//TODO: fill out
func (i *Interface) Attrs() gographviz.Attrs {
	return nil
}

func (i *Interface) Edges() []*gographviz.Edge {
	var retval []*gographviz.Edge

	retval = make([]*gographviz.Edge, 0, len(i.inheritedInterfaces)+len(i.implementedByCache))

	for _, v := range i.inheritedInterfaces {
		//TODO: decide on attrs
		retval = append(retval, &gographviz.Edge{v.Name(), "", i.Name(), "", true, nil})
	}
	for _, v := range i.implementedByCache {
		//TODO: decide on attrs
		retval = append(retval, &gographviz.Edge{v.Name(), "", i.Name(), "", true, nil})
	}

	return retval
}

func (i *Interface) isComposite() bool {
	return len(i.inheritedInterfaces) > 0
}

//no mutation
func (i *Interface) isImplementedBy(s *Struct) bool {
	required := i.allRequiredFunctions()
	have := s.allReceiverFunctions()

	fmt.Println("struct:", s.Name())
	fmt.Println("interface:", i.Name())
	fmt.Println("required:", required)
	fmt.Println("have:", have)

	for _, v := range required {
		found := false
		for _, j := range have {
			if j == v {
				found = true
				break
			}
		}

		if !found {
			fmt.Println("return false\n")
			return false
		}
	}

	fmt.Println("return true\n")
	return true
}

func (i *Interface) setImplementedBy(s []*Struct) []*Struct {
	retval := i.implementedBy(s)

	i.implementedByCache = make([]*Struct, len(retval))
	copy(i.implementedByCache, retval)

	return retval
}

//filter pattern
func (i *Interface) implementedBy(s []*Struct) []*Struct {
	var retval []*Struct

	for _, v := range s {
		if i.isImplementedBy(v) {
			retval = append(retval, v)
		}
	}

	return retval
}

//no mutation
func (i *Interface) allRequiredFunctions() []*Function {
	retval := make([]*Function, len(i.requiredFunctions))
	c := copy(retval, i.requiredFunctions)
	if c != len(i.requiredFunctions) {
		panic("copy failed in allRequiredFunctions")
	}

	for _, v := range i.inheritedInterfaces {
		retval = append(retval, v.allRequiredFunctions()...)
	}

	return retval
}

//for if interface is found as an Anonymous member of something else first
func makeInterfaceUnknown(source *Interface, b *BaseType) *Interface {
	retval := &Interface{b, make([]*Function, 0), make([]*Interface, 0), make([]gographviz.GraphableNode, 0), nil, nil, nil}
	b.addNode(retval)

	retval.addToIncludedIn(source)

	return retval
}

func (i *Interface) addToIncludedIn(x gographviz.GraphableNode) {
	i.includedIn = append(i.includedIn, x)
}

func (i *Interface) remakeInterface(spec *ast.TypeSpec) *Interface {
	interfaceType, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		panic("bad ast.TypeSpec that is not InterfaceType in makeInterface")
	}

	for _, v := range interfaceType.Methods.List {
		if len(v.Names) != 0 {
			f := funcMap.lookupOrAddFromExpr(v.Names[0].Name, v.Type.(*ast.FuncType))
			f.addInterface(i)
			i.requiredFunctions = append(i.requiredFunctions, f)
		} else {
			lookup := typeMap.lookupOrAdd(String(v.Type))
			node := lookup.base.node
			if node != nil {
				i.inheritedInterfaces = append(i.inheritedInterfaces, node.(*Interface))
				node.(*Interface).addToIncludedIn(i)
			} else {
				i.inheritedInterfaces = append(i.inheritedInterfaces, makeInterfaceUnknown(i, lookup.base))
			}
		}
	}

	return i
}

//possibilities for lines:
//Type -> inheritedStruct
//(comma seperated list of names) Type -> NamedTypes
//b: the baseType for this struct
//lines: lines from the structs declaration block, preceeding and trailing whitespace removed
func makeInterface(spec *ast.TypeSpec, b *BaseType) *Interface {
	interfaceType, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		panic("bad ast.TypeSpec that is not InterfaceType in makeInterface")
	}

	//should only be used with declarations, if struct is in field names use makeStructUnknown
	retval := &Interface{b, make([]*Function, 0), make([]*Interface, 0), make([]gographviz.GraphableNode, 0), nil, nil, interfaceType}

	for _, v := range interfaceType.Methods.List {
		if len(v.Names) != 0 {
			f := funcMap.lookupOrAddFromExpr(v.Names[0].Name, v.Type.(*ast.FuncType))
			f.addInterface(retval)
			retval.requiredFunctions = append(retval.requiredFunctions, f)
		} else {
			lookup := typeMap.lookupOrAdd(String(v.Type))
			node := lookup.base.node
			if node != nil {
				retval.inheritedInterfaces = append(retval.inheritedInterfaces, node.(*Interface))
				node.(*Interface).addToIncludedIn(retval)
			} else {
				retval.inheritedInterfaces = append(retval.inheritedInterfaces, makeInterfaceUnknown(retval, lookup.base))
			}
		}
	}

	//fmt.Println("makeInterface:")
	//fmt.Println(retval)
	b.addNode(retval)
	return retval
}
