package main

type Struct struct {
	target Type

	fields []NameTypePair

	//receiver functions that only work with pointer to this type
	pointerReceiverFunctions []ReceiverFunction

	//functions that 
	valueReceiverFunctions []ReceiverFunction
}