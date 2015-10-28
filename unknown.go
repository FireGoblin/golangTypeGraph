package main

import (
	"github.com/firegoblin/gographviz"
	"go/ast"
)

//A node is responsible for the incoming edges to it

//A node type
//implements gographviz.GraphableNode
type unknownNode struct {
	target *baseType
}

func (u *unknownNode) Name() string {
	return gographviz.SafeName(u.target.name)
}

func (u *unknownNode) Attrs() gographviz.Attrs {
	return nil
}

func (u *unknownNode) Edges() []*gographviz.Edge {
	return nil
}

//for if struct is found as an Anonymous member of something else first
func newUnknown(source *structNode, b *baseType) *unknownNode {
	retval := &unknownNode{b}
	b.addNode(retval)

	return retval
}

func (u *unknownNode) remakeStruct(spec *ast.TypeSpec) *structNode {
	retval := &structNode{u.target, nil, make([]namedType, 0), make([]receiverFunction, 0), make([]*baseType, 0), nil, nil, spec.Type}

	retval.remakeStructInternals(spec)

	retval.target.addNode(retval)
	return retval
}

func (u *unknownNode) remakeInterface(spec *ast.TypeSpec) *interfaceNode {
	interfaceType, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		panic("bad ast.TypeSpec that is not InterfaceType in newInterface")
	}

	retval := &interfaceNode{u.target, make([]*function, 0), make([]*interfaceNode, 0), nil, nil, interfaceType}

	retval.remakeInterfaceInternals(interfaceType)

	retval.target.addNode(retval)
	return retval
}
