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

	astNode ast.FuncType
}

func makeFunction(s string, f *ast.FuncType) *Function {
	//FunctionParser := MustCompile(`^([\w]+)(\(.*?\) .*)$`)

	//reg := FunctionParser.FindStringSubmatch(s)
	typ := typeMap.lookupOrAdd("func" + String(f))
	params, err := typ.params()
	if err != nil {
		panic(err)
	}
	returns, err := typ.returnTypes()
	if err != nil {
		panic(err)
	}

	retval := Function{s, typ, params, returns, make([]*Struct, 0), make([]*Interface, 0), *f}

	return &retval
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
