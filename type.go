package main

import "go/ast"

type Type struct {
	name    string
	base    *BaseType
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
	var baseType string
	if expr != nil {
		baseType = String(BaseTypeOf(expr))
	} else {
		baseType = s
	}

	retval := Type{s, nil, expr}

	if RecursiveTypeOf(expr) == nil {
		b := makeBase(baseType)
		retval.base = b
	} else {
		b, ok := typeMap[baseType]
		if !ok {
			b = typeMap.lookupOrAdd(baseType)
		}

		retval.base = b.base

		next := RecursiveTypeOf(expr)
		//create lower type if not created yet
		_, ok = typeMap[String(next)]
		if !ok {
			makeTypeRecursive(String(next), retval.base, next)
		}
	}

	return &retval
}

//never call outside of makeType
func makeTypeRecursive(s string, b *BaseType, expr ast.Expr) {
	x := Type{s, b, expr}
	typeMap[s] = &x

	next := RecursiveTypeOf(expr)
	if next != nil {
		_, ok := typeMap[String(next)]
		if !ok {
			makeTypeRecursive(String(next), b, next)
		}
	}
}

func (t Type) String() string {
	return t.name
}
