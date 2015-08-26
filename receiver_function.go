package main 

//represents a receiver function
//A single copy should exist no matter how many structs and interfaces use it
type ReceiverFunction struct {
	InterfaceFunction 

	//all structs that implement this function
	structs []Struct

	//any interfaces that require this function
	interfaces []Interface 
}