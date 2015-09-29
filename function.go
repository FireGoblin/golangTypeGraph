package main

import . "regexp"

var FunctionParser = MustCompile(`^([\w]+)(\(.*?\) .*)$`)

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
}

func makeFunction(s string) *Function {
	reg := FunctionParser.FindStringSubmatch(s)
	typ := typeMap.lookupOrAdd("func" + reg[2])
	params, err := typ.params()
	if err != nil {
		panic("Type for function being made is not a Function type")
	}
	returns, _ := typ.returnTypes()

	retval := Function{reg[1], typ, params, returns, make([]*Struct, 0), make([]*Interface, 0)}

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
