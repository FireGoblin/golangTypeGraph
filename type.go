package main

import "go/ast"

type Type struct {
	name    string
	base    *BaseType
	astNode ast.Expr
}

//TODO: maybe more effecient way to do this
//only call from master_type_map
func makeType(s string, pkg string) *Type {
	return sharedMakeType(s, nil, pkg)
}

//TODO: maybe more effecient way to do this
//only call from master_type_map
func makeTypeFromExpr(expr ast.Expr, pkg string) *Type {
	return sharedMakeType(String(expr), expr, pkg)
}

func sharedMakeType(s string, expr ast.Expr, pkg string) *Type {
	var baseType string
	if expr != nil {
		baseType = String(BaseTypeOf(expr))
	} else {
		baseType = s
	}

	retval := Type{s, nil, expr}

	if RecursiveTypeOf(expr) == nil {
		b := makeBase(baseType, pkg)
		retval.base = b
	} else {
		b, ok := typeMap.getPkg(pkg)[baseType]
		if !ok {
			b = typeMap.lookupOrAddWithPkg(baseType, pkg)
		}

		retval.base = b.base

		next := RecursiveTypeOf(expr)
		//create lower type if not created yet
		_, ok = typeMap.getPkg(pkg)[String(next)]
		if !ok {
			makeTypeRecursive(String(next), retval.base, next, pkg)
		}
	}

	return &retval
}

//never call outside of makeType
func makeTypeRecursive(s string, b *BaseType, expr ast.Expr, pkg string) {
	x := Type{s, b, expr}
	typeMap.theMap[pkg][s] = &x

	next := RecursiveTypeOf(expr)
	if next != nil {
		_, ok := typeMap.theMap[pkg][String(next)]
		if !ok {
			makeTypeRecursive(String(next), b, next, pkg)
		}
	}
}

func (t Type) String() string {
	return t.name
}
