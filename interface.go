package main 

type Interface struct {
	target Type 

	//any functions required to implement the interface
	requiredFunctions []ReceiverFunction

	//functions that act on this type
	ownFunctions []ReceiverFunction

	//interfaces this inherits from
	//if not zero, is composite interface
	inheritedInterfaces []Interface
}