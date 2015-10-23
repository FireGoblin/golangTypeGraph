package main

//import . "regexp"
import "strings"
import "go/ast"

type Type struct {
	name string
	base *BaseType

	//corresponds to the number of asterisks for the type
	//exp: **string would have pointerLevel = 2
	pointerLevel int

	astNode ast.Expr
}

//TODO: maybe more effecient way to do this
//only call from master_type_map
func makeType(s string) *Type {
	return sharedMakeType(s, nil)
}

//TODO: maybe more effecient way to do this
//only call from master_type_map
func makeTypeFromExpr(expr ast.Expr) *Type {
	return sharedMakeType(String(expr), expr)
}

func sharedMakeType(s string, expr ast.Expr) *Type {
	baseType := strings.Trim(s, "*")
	pLevel := len(s) - len(baseType)

	retval := Type{s, nil, pLevel, expr}

	if pLevel == 0 {
		b := makeBase(baseType)
		retval.base = b
		b.addType(&retval)
	} else {
		b, ok := typeMap[baseType]
		if !ok {
			b = typeMap.lookupOrAdd(baseType)
		}

		retval.base = b.base
		b.base.addType(&retval)

		//create lower type if not created yet
		_, ok = typeMap[s[1:]]
		if !ok {
			makeTypeRecursive(s[1:], retval.base, pLevel-1, expr.(*ast.StarExpr).X)
		}
	}

	return &retval
}

//never call outside of makeType
func makeTypeRecursive(s string, b *BaseType, pLevel int, expr ast.Expr) {
	x := Type{s, b, pLevel, expr}
	typeMap[s] = &x
	b.addType(&x)

	_, ok := typeMap[s[1:]]
	if !ok {
		makeTypeRecursive(s[1:], b, pLevel-1, expr.(*ast.StarExpr).X)
	}
}

func (t Type) String() string {
	return t.name
}
