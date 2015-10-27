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

	s := StringInterfaceField(name, namelessExpr)

	x, ok := m[s]

	if ok && x.astNode == nil {
		x.astNode = expr
	} else if !ok {
		m[s] = newFunction(name, expr, namelessExpr)

		//error checking
		x, ok = m[s]
		if !ok || x == nil {
			panic("masterlist not properly associated with new func")
		}
	}

	return x
}
