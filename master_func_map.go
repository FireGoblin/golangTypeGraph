package main

import "go/ast"

//master map uses singleton pattern, only one of them should be created in the program
//only master should call creators for types
type masterFuncMap map[string]*Function

var funcMap = masterFuncMap(make(map[string]*Function))

func (m masterFuncMap) lookupOrAddFromExpr(name string, expr *ast.FuncType) *Function {
	namelessExpr := &ast.FuncType{0, nil, nil}

	namelessExpr.Params = Normalized(expr.Params)
	namelessExpr.Results = Normalized(expr.Results)

	newFunc := newFunction(name, expr, namelessExpr)

	s := newFunc.lookupString()

	x, ok := m[s]

	if ok && x.astNode == nil {
		x.astNode = expr
	} else if !ok {
		m[s] = newFunc

		//error checking
		x, ok = m[s]
		if !ok || x == nil {
			panic("masterlist not properly associated with new func")
		}
	}

	return x
}
