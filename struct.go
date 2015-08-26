package main

//A node type
type Struct struct {
	target BaseType

	fields []NameTypePair

	//receiver functions that only work with pointer to this type
	pointerReceiverFunctions []ReceiverFunction

	//functions that 
	valueReceiverFunctions []ReceiverFunction
}