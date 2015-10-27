package main

import "go/ast"

type Function struct {
	//name for function
	name string

	//type of the function
	//should pass target.isFunc() == true, otherwise panic
	target *Type

	//The types of the params to this function
	paramTypes []*Type

	//The types of the return value from this function
	returnTypes []*Type

	isReceiver bool

	astNode *ast.FuncType
}

func (f *Function) String() string {
	return f.name + " " + f.target.name
}

func makeFunction(s string, f *ast.FuncType, nameless *ast.FuncType) *Function {
	typ := typeMap.lookupOrAddFromExpr(f)

	var paramsProcessed = make([]*Type, 0)
	var resultsProcessed = make([]*Type, 0)

	if nameless.Params != nil {
		for _, expr := range nameless.Params.List {
			paramsProcessed = append(paramsProcessed, typeMap.lookupOrAddFromExpr(expr.Type))
		}
	}

	if nameless.Results != nil {
		for _, expr := range nameless.Results.List {
			resultsProcessed = append(resultsProcessed, typeMap.lookupOrAddFromExpr(expr.Type))
		}
	}

	//TODO eventually: re-add paramTypes and returnTypes
	retval := &Function{s, typ, paramsProcessed, resultsProcessed, false, f}

	return retval
}
