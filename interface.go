package main

import (
	"github.com/firegoblin/gographviz"
	"go/ast"
)

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

	//structs
	//Draw edges from these
	implementedByCache []*Struct

	extraAttrs gographviz.Attrs

	astNode *ast.InterfaceType
}

func (i *Interface) String() string {
	return i.target.String()
}

func (i *Interface) Name() string {
	return i.target.Name()
}

func (i *Interface) Attrs() gographviz.Attrs {
	retval := make(map[string]string)
	retval["shape"] = "Mrecord"
	retval["label"] = i.label()
	return retval
}

func (i *Interface) Edges() []*gographviz.Edge {
	var retval []*gographviz.Edge

	retval = make([]*gographviz.Edge, 0, len(i.inheritedInterfaces)+len(i.implementedByCache))

	for _, v := range i.inheritedInterfaces {
		//TODO: decide on attrs
		retval = append(retval, &gographviz.Edge{v.Name(), "", i.Name(), "", true, inheritedAttrs()})
	}
	for _, v := range i.implementedByCache {
		//TODO: decide on attrs
		retval = append(retval, &gographviz.Edge{v.Name(), "", i.Name(), "", true, i.implementedAttrs()})
	}

	return retval
}

func (i *Interface) highlyImplemented() bool {
	return len(i.implementedByCache) > *implementMax
}

func (i *Interface) label() string {
	retval := "\"{" + i.String() + "|"

	if i.highlyImplemented() {
		retval += "*HIGHLY IMPLMENTED*\\n"
	}

	for _, v := range i.inheritedInterfaces {
		retval += v.target.stringRelativePkg(i.target.pkgName) + "\\n"
	}

	retval += "|"

	for _, v := range i.requiredFunctions {
		retval += v.String() + "\\l"
	}

	retval += "}\""

	return retval
}

func (i *Interface) implementedAttrs() map[string]string {
	retval := make(map[string]string)
	retval["label"] = "implements"
	if i.highlyImplemented() {
		retval["style"] = "invis"
	} else {
		retval["style"] = "bold"
	}
	return retval
}

//no mutation
func (i *Interface) isImplementedBy(s *Struct) bool {
	required := i.allRequiredFunctions()
	have := s.allreceiverFunctions()

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
func newInterfaceUnknown(source *Interface, b *BaseType) *Interface {
	retval := &Interface{b, make([]*Function, 0), make([]*Interface, 0), nil, nil, nil}
	b.addNode(retval)

	return retval
}

func (i *Interface) remakeInterfaceInternals(interfaceType *ast.InterfaceType) {
	for _, v := range interfaceType.Methods.List {
		if len(v.Names) != 0 {
			f := funcMap.lookupOrAddFromExpr(v.Names[0].Name, v.Type.(*ast.FuncType))
			f.isReceiver = true
			i.requiredFunctions = append(i.requiredFunctions, f)
		} else {
			lookup := typeMap.lookupOrAddFromExpr(v.Type)
			node := lookup.base.node
			if node != nil {
				i.inheritedInterfaces = append(i.inheritedInterfaces, node.(*Interface))
			} else {
				i.inheritedInterfaces = append(i.inheritedInterfaces, newInterfaceUnknown(i, lookup.base))
			}
		}
	}
}

func (i *Interface) remakeInterface(spec *ast.TypeSpec) *Interface {
	interfaceType, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		panic("bad ast.TypeSpec that is not InterfaceType in newInterface")
	}

	i.remakeInterfaceInternals(interfaceType)

	return i
}

//possibilities for lines:
//Type -> inheritedStruct
//(comma seperated list of names) Type -> namedTypes
//b: the baseType for this struct
//lines: lines from the structs declaration block, preceeding and trailing whitespace removed
func newInterface(spec *ast.TypeSpec, b *BaseType) *Interface {
	interfaceType, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		panic("bad ast.TypeSpec that is not InterfaceType in newInterface")
	}

	//should only be used with declarations, if struct is in field names use newStructUnknown
	retval := &Interface{b, make([]*Function, 0), make([]*Interface, 0), nil, nil, interfaceType}

	retval.remakeInterfaceInternals(interfaceType)

	b.addNode(retval)
	return retval
}
