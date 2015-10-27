package main

import "github.com/firegoblin/gographviz"
import "go/ast"
import "fmt"

//data type for fields/parameters.
//A pair of the name of the variable and its type.
type NamedType struct {
	name   string
	target *Type
}

func (n NamedType) String() string {
	return n.name + " " + n.target.String()
}

func (n NamedType) StringRelativePkg(pkg string) string {
	return n.name + " " + n.target.StringRelativePkg(pkg)
}

func (n NamedType) Node() gographviz.GraphableNode {
	return n.target.base.node
}

func NamedTypeFromField(f *ast.Field) NamedType {
	if len(f.Names) != 1 {
		panic(fmt.Sprintf("tried to created NamedType with %d names", len(f.Names)))
	}

	return NamedType{f.Names[0].Name, typeMap.lookupOrAddFromExpr(f.Type)}
}
