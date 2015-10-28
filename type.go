package main

import "go/ast"

type Type struct {
	name    string
	base    *baseType
	astNode ast.Expr
}

//TODO: maybe more effecient way to do this
//only call from master_type_map
func newType(s string, pkg string) *Type {
	return sharedNewType(s, nil, pkg)
}

//TODO: maybe more effecient way to do this
//only call from master_type_map
func newTypeFromExpr(expr ast.Expr, pkg string) *Type {
	return sharedNewType(String(expr), expr, pkg)
}

func sharedNewType(s string, expr ast.Expr, pkg string) *Type {
	var baseType string
	if expr != nil {
		baseType = String(RootTypeOf(expr))
	} else {
		baseType = s
	}

	retval := Type{s, nil, expr}

	if RecursiveTypeOf(expr) == nil {
		b := newBase(baseType, pkg)
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
			newTypeRecursive(String(next), retval.base, next, pkg)
		}
	}

	return &retval
}

//never call outside of newType
func newTypeRecursive(s string, b *baseType, expr ast.Expr, pkg string) {
	x := Type{s, b, expr}
	typeMap.theMap[pkg][s] = &x

	next := RecursiveTypeOf(expr)
	if next != nil {
		_, ok := typeMap.theMap[pkg][String(next)]
		if !ok {
			newTypeRecursive(String(next), b, next, pkg)
		}
	}
}

func (t Type) String() string {
	return t.stringRelativePkg(*defaultPackageName)
}

func (t Type) stringRelativePkg(pkg string) string {
	if pkg == t.base.pkgName {
		return t.name
	}

	return stringWithPkg(t.base.pkgName, t.astNode)
}
