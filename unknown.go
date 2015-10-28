package main

import (
	"github.com/firegoblin/gographviz"
	"go/ast"
)

//A node is responsible for the incoming edges to it

//A node type
//implements gographviz.GraphableNode
type unknown struct {
	target *baseType
}

func (u *unknown) Name() string {
	return gographviz.SafeName(u.target.name)
}

func (u *unknown) Attrs() gographviz.Attrs {
	return nil
}

func (u *unknown) Edges() []*gographviz.Edge {
	return nil
}

//for if struct is found as an Anonymous member of something else first
func newUnknown(source *Struct, b *baseType) *unknown {
	retval := &unknown{b}
	b.addNode(retval)

	return retval
}

func (u *unknown) remakeStruct(spec *ast.TypeSpec) *Struct {
	retval := &Struct{u.target, nil, make([]namedType, 0), make([]receiverFunction, 0), make([]*baseType, 0), nil, nil, spec.Type}

	retval.remakeStructInternals(spec)

	retval.target.addNode(retval)
	return retval
}

func (u *unknown) remakeInterface(spec *ast.TypeSpec) *Interface {
	interfaceType, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		panic("bad ast.TypeSpec that is not InterfaceType in newInterface")
	}

	retval := &Interface{u.target, make([]*function, 0), make([]*Interface, 0), nil, nil, interfaceType}

	retval.remakeInterfaceInternals(interfaceType)

	retval.target.addNode(retval)
	return retval
}
