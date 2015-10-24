package main

import "go/ast"

//import "fmt"

//master map uses singleton pattern, only one of them should be created in the program
//only master should call creators for types
type MasterFuncMap map[string]*Function

var funcMap = MasterFuncMap(make(map[string]*Function))

func (m MasterFuncMap) lookupOrAddFromExpr(name string, expr *ast.FuncType) *Function {
	var namelessExpr = &ast.FuncType{0, nil, nil}

	namelessExpr.Params = normalized(expr.Params)
	namelessExpr.Results = normalized(expr.Results)

	s := StringInterfaceField(name, namelessExpr)

	x, ok := m[s]

	if ok && x.astNode == nil {
		x.astNode = expr
	} else if !ok {
		m[s] = makeFunction(name, expr, namelessExpr)

		//error checking
		x, ok = m[s]
		if !ok || x == nil {
			panic("masterlist not properly associated with new func")
		}
	}

	return x
}
