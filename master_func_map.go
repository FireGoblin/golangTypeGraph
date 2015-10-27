package main

import "go/ast"

//master map uses singleton pattern, only one of them should be created in the program
//only master should call creators for types
type masterFuncMap struct {
	theMap     map[string]map[string]*Function
	currentPkg string
}

var funcMap = masterFuncMap{make(map[string]map[string]*Function), ""}

func (m masterFuncMap) currentMap() map[string]*Function {
	_, ok := m.theMap[m.currentPkg]
	if !ok {
		m.theMap[m.currentPkg] = make(map[string]*Function)
	}
	return m.theMap[m.currentPkg]
}

func (m masterFuncMap) getPkg(pkg string) map[string]*Function {
	_, ok := m.theMap[pkg]
	if !ok {
		m.theMap[pkg] = make(map[string]*Function)
	}
	return m.theMap[pkg]
}

func (m masterFuncMap) lookupOrAddFromExpr(name string, expr *ast.FuncType) *Function {
	namelessExpr := &ast.FuncType{0, nil, nil}

	namelessExpr.Params = Normalized(expr.Params)
	namelessExpr.Results = Normalized(expr.Results)

	s := StringInterfaceField(name, namelessExpr)

	x, ok := m.currentMap()[s]

	if ok && x.astNode == nil {
		x.astNode = expr
	} else if !ok {
		m.currentMap()[s] = newFunction(name, expr, namelessExpr)

		//error checking
		x, ok = m.currentMap()[s]
		if !ok || x == nil {
			panic("masterlist not properly associated with new func")
		}
	}

	return x
}
