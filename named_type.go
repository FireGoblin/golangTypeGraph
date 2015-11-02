package main

import (
	"fmt"
	"github.com/firegoblin/gographviz"
	"go/ast"
	"go/token"
	"strings"
)

//data type for fields/parameters.
//A pair of the name of the variable and its type.
type namedType struct {
	name   string
	target *Type
	tag    string
}

func (n namedType) String() string {
	return n.stringShared(n.target.String())
}

func (n namedType) stringRelativePkg(pkg string) string {
	return n.stringShared(n.target.stringRelativePkg(pkg))
}

func (n namedType) stringShared(targetStr string) string {
	if n.name == "" {
		return targetStr
	}

	str := n.name + " " + targetStr
	if n.tag != "" && *json {
		str += " " + n.tag
	}
	return str
}

func (n namedType) node() gographviz.GraphableNode {
	return n.target.base.node
}

func stringFromBasicLit(basic *ast.BasicLit) string {
	if basic == nil || basic.Kind != token.STRING {
		return ""
	}

	return strings.Replace(basic.Value, "\"", "\\\"", -1)
}

func newNamedTypeFromField(f *ast.Field) namedType {
	if len(f.Names) > 1 {
		panic(fmt.Sprintf("tried to create namedType with %d names", len(f.Names)))
	} else if len(f.Names) == 0 {
		return namedType{"", typeMap.lookupOrAddFromExpr(f.Type), stringFromBasicLit(f.Tag)}
	}

	return namedType{f.Names[0].Name, typeMap.lookupOrAddFromExpr(f.Type), stringFromBasicLit(f.Tag)}
}
