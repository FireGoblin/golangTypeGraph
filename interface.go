package main

import "go/ast"
import "github.com/firegoblin/gographviz"

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
	includedIn []*Interface

	//structs
	//Draw edges from these
	implementedByCache []*Struct

	extraAttrs gographviz.Attrs

	astNode *ast.InterfaceType
}

func (i *Interface) String() string {
	return i.target.name
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

func (i *Interface) setImplementedBy(s []*Struct) []*Struct {
	retval := i.implementedBy(s)

	i.implementedByCache = make([]*Struct, 0, len(retval))
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
func makeInterfaceUnknown(b *BaseType, source *Interface) *Interface {
	retval := Interface{b, make([]*Function, 0), make([]*Interface, 0), make([]*Interface, 0), nil, nil, nil}
	b.node = &retval

	retval.includedIn = append(retval.includedIn, source)

	return &retval
}

//possibilities for lines:
//Type -> inheritedStruct
//(comma seperated list of names) Type -> NamedTypes
//b: the baseType for this struct
//lines: lines from the structs declaration block, preceeding and trailing whitespace removed
func makeInterface(spec *ast.TypeSpec, typ *BaseType) *Interface {
	// FunctionParser := MustCompile(`^([\w]+)(\(.*?\) .*)$`)
	// retval := Interface{b, make([]*Function, 0), make([]*Interface, 0), make([]*Interface, 0)}
	// b.node = &retval

	// for _, v := range lines {
	// 	ifp := FunctionParser.FindStringSubmatch(v)
	// 	if len(ifp) != 0 {
	// 		f := funcMap.lookupOrAdd(ifp[0])
	// 		f.addInterface(&retval)
	// 		retval.requiredFunctions = append(retval.requiredFunctions, f)
	// 	} else {
	// 		typ := typeMap.lookupOrAdd(v)
	// 		var interfac *Interface

	// 		if typ.base.node == nil {
	// 			interfac := makeInterfaceUnknown(typ.base, &retval)
	// 			typ.base.node = interfac
	// 		} else {
	// 			var ok bool
	// 			interfac, ok = typ.base.node.(*Interface)
	// 			if !ok {
	// 				panic("Could not find struct of anonymous member")
	// 			}
	// 		}
	// 		retval.inheritedInterfaces = append(retval.inheritedInterfaces, interfac)
	// 	}
	// }

	// return &retval

	return nil
}
