package main

import (
	"fmt"
	"github.com/firegoblin/gographviz"
	"go/ast"
)

//data type for fields/parameters.
//A pair of the name of the variable and its type.
type namedType struct {
	name   string
	target *Type
}

func (n namedType) String() string {
	return n.name + " " + n.target.String()
}

func (n namedType) StringRelativePkg(pkg string) string {
	return n.name + " " + n.target.StringRelativePkg(pkg)
}

func (n namedType) Node() gographviz.GraphableNode {
	return n.target.base.node
}

func newNamedTypeFromField(f *ast.Field) namedType {
	if len(f.Names) != 1 {
		panic(fmt.Sprintf("tried to created namedType with %d names", len(f.Names)))
	}

	return namedType{f.Names[0].Name, typeMap.lookupOrAddFromExpr(f.Type)}
}
