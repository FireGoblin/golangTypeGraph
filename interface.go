package main 

//A node type
type Interface struct {
	target BaseType 

	//any functions required to implement the interface
	requiredFunctions []ReceiverFunction

	//interfaces this inherits from
	//if not zero, is composite interface
	inheritedInterfaces []Interface
}