package main

type InterfaceFunction struct {
	//name for function
	name string

	//type of the function
	target Type

	paramTypes []Type 
	returnTypes []Type
}