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

	//added to handle files where receiver methods declared on a struct not in the file
	receiverFunctions []receiverFunction
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
	retval := &unknownNode{b, make([]receiverFunction, 0)}
	b.addNode(retval)

	return retval
}

func (u *unknownNode) addFunction(f *function, field *ast.Field) {
	u.receiverFunctions = append(u.receiverFunctions, newReceiverFunction(f, field))
	f.isReceiver = true
}

func (u *unknownNode) remakeStruct(spec *ast.TypeSpec) *structNode {
	retval := &structNode{u.target, nil, make([]namedType, 0), make([]receiverFunction, len(u.receiverFunctions)), make([]*baseType, 0), nil, nil, spec.Type}

	retval.remakeStructInternals(spec)

	copy(retval.receiverFunctions, u.receiverFunctions)
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
