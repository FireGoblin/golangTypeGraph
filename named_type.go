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
	if n.name == "" {
		return n.target.String()
	}

	return n.name + " " + n.target.String()
}

func (n namedType) stringRelativePkg(pkg string) string {
	return n.name + " " + n.target.stringRelativePkg(pkg)
}

func (n namedType) node() gographviz.GraphableNode {
	return n.target.base.node
}

func newNamedTypeFromField(f *ast.Field) namedType {
	if len(f.Names) > 1 {
		panic(fmt.Sprintf("tried to create namedType with %d names", len(f.Names)))
	} else if len(f.Names) == 0 {
		return namedType{"", typeMap.lookupOrAddFromExpr(f.Type)}
	}

	return namedType{f.Names[0].Name, typeMap.lookupOrAddFromExpr(f.Type)}
}
