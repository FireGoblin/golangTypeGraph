package main

import (
	"github.com/firegoblin/gographviz"
	"go/ast"
)

//A node type
//implements gographviz.GraphableNode
type interfaceNode struct {
	target *baseType

	//any functions required to implement the interface
	//does not include functions inherited indirectly through inheritedInterfaces
	requiredFunctions []*function

	//interfaces this inherits from
	//if not zero, is composite interface
	//Draw edges from these
	inheritedInterfaces []*interfaceNode

	//structs
	//Draw edges from these
	implementedByCache []*structNode

	extraAttrs gographviz.Attrs

	astNode *ast.InterfaceType
}

func (i *interfaceNode) String() string {
	return i.target.String()
}

func (i *interfaceNode) Name() string {
	return i.target.Name()
}

func (i *interfaceNode) Attrs() gographviz.Attrs {
	retval := make(map[string]string)
	retval["shape"] = "Mrecord"
	retval["label"] = i.label()
	return retval
}

func (i *interfaceNode) Edges() []*gographviz.Edge {
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

func (i *interfaceNode) highlyImplemented() bool {
	return len(i.implementedByCache) > *implementMax
}

func (i *interfaceNode) label() string {
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

func (i *interfaceNode) implementedAttrs() map[string]string {
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
func (i *interfaceNode) isImplementedBy(s *structNode) bool {
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

func (i *interfaceNode) setImplementedBy(s []*structNode) []*structNode {
	retval := i.implementedBy(s)

	i.implementedByCache = make([]*structNode, len(retval))
	copy(i.implementedByCache, retval)

	return retval
}

//filter pattern
func (i *interfaceNode) implementedBy(s []*structNode) []*structNode {
	var retval []*structNode

	for _, v := range s {
		if i.isImplementedBy(v) {
			retval = append(retval, v)
		}
	}

	return retval
}

//no mutation
func (i *interfaceNode) allRequiredFunctions() []*function {
	retval := make([]*function, len(i.requiredFunctions))
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
func newInterfaceUnknown(source *interfaceNode, b *baseType) *interfaceNode {
	retval := &interfaceNode{b, make([]*function, 0), make([]*interfaceNode, 0), nil, nil, nil}
	b.addNode(retval)

	return retval
}

func (i *interfaceNode) remakeInterfaceInternals(interfaceType *ast.InterfaceType) {
	for _, v := range interfaceType.Methods.List {
		if len(v.Names) != 0 {
			f := funcMap.lookupOrAddFromExpr(v.Names[0].Name, v.Type.(*ast.FuncType))
			f.isReceiver = true
			i.requiredFunctions = append(i.requiredFunctions, f)
		} else {
			lookup := typeMap.lookupOrAddFromExpr(v.Type)
			node := lookup.base.node
			if node != nil {
				i.inheritedInterfaces = append(i.inheritedInterfaces, node.(*interfaceNode))
			} else {
				i.inheritedInterfaces = append(i.inheritedInterfaces, newInterfaceUnknown(i, lookup.base))
			}
		}
	}
}

func (i *interfaceNode) remakeInterface(spec *ast.TypeSpec) *interfaceNode {
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
func newInterface(spec *ast.TypeSpec, b *baseType) *interfaceNode {
	interfaceType, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		panic("bad ast.TypeSpec that is not InterfaceType in newInterface")
	}

	//should only be used with declarations, if struct is in field names use newStructUnknown
	retval := &interfaceNode{b, make([]*function, 0), make([]*interfaceNode, 0), nil, nil, interfaceType}

	retval.remakeInterfaceInternals(interfaceType)

	b.addNode(retval)
	return retval
}
