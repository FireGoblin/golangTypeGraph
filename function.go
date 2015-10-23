package main

//import . "regexp"
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

	//if interface or struct is non-empty, then this is used as a receiver function
	//all structs that implement this function
	structs []*Struct

	//any interfaces that require this function
	interfaces []*Interface

	astNode *ast.FuncType
}

func makeFunction(s string, f *ast.FuncType, nameless *ast.FuncType) *Function {
	typ := typeMap.lookupOrAdd(String(f))

	var paramsProcessed = make([]*Type, 0)
	var resultsProcessed = make([]*Type, 0)

	for _, expr := range nameless.Params.List {
		paramsProcessed = append(paramsProcessed, typeMap.lookupOrAddFromExpr(expr.Type))
	}
	for _, expr := range nameless.Results.List {
		resultsProcessed = append(resultsProcessed, typeMap.lookupOrAddFromExpr(expr.Type))
	}

	//TODO eventually: re-add paramTypes and returnTypes
	retval := &Function{s, typ, paramsProcessed, resultsProcessed, make([]*Struct, 0), make([]*Interface, 0), f}

	return retval
}

func (f *Function) addInterface(i *Interface) {
	f.interfaces = append(f.interfaces, i)
}

func (f *Function) addStruct(s *Struct) {
	f.structs = append(f.structs, s)
}

func (f Function) isReceiverFunction() bool {
	return len(f.structs)+len(f.interfaces) > 0
}
