package main

import "go/ast"

//master map uses singleton pattern, only one of them should be created in the program
//only master should call creators for types
type MasterFuncMap map[string]*Function

var funcMap = MasterFuncMap(make(map[string]*Function))

// func (m MasterFuncMap) lookupOrAdd(s string) *Function {
// 	x, ok := m[s]

// 	if !ok {
// 		m[s] = makeFunction(s)

// 		//error checking
// 		x, ok = m[s]
// 		if !ok || x == nil {
// 			panic("masterlist not properly associated with new func")
// 		}
// 	}

// 	return x
// }

func (m MasterFuncMap) lookupOrAddFromExpr(name string, expr *ast.FuncType) *Function {
	s := name + String(expr)

	x, ok := m[s]

	if ok && x.astNode == nil {
		x.astNode = expr
	} else {
		m[s] = makeFunction(name, expr)

		//error checking
		x, ok = m[s]
		if !ok || x == nil {
			panic("masterlist not properly associated with new func")
		}
	}

	return x
}
