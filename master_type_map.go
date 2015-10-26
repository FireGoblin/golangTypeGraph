package main

import "go/ast"

//master map uses singleton pattern, only one of them should be created in the program
//only master should call creators for types
type MasterTypeMap map[string]*Type

var typeMap = MasterTypeMap(make(map[string]*Type))

//TODO: Consider whether to attach the object
//only call if the string is known to be a base type
func (m MasterTypeMap) lookupOrAdd(s string) *Type {
	x, ok := m[s]

	if !ok {
		m[s] = makeType(s)

		x, ok = m[s]
		if !ok {
			panic("masterlist not properly associated with new type")
		}
	}

	return x
}

func (m MasterTypeMap) lookupOrAddFromExpr(expr ast.Expr) *Type {
	s := String(expr)

	x, ok := m[s]

	if ok && x.astNode == nil {
		x.astNode = expr
	} else if !ok {
		m[s] = makeTypeFromExpr(expr)

		//error checking
		x, ok = m[s]
		if !ok || x == nil {
			panic("masterlist not properly associated with new func")
		}
	}

	return x
}
