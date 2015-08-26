package main

import . "regexp"

var AnonymousStructMatcher = MustCompile(`^[^ ]$`)

//A node type
type Struct struct {
	target *BaseType

	fields []NameTypePair

	//receiver functions that only work with pointer to this type
	pointerReceiverFunctions []ReceiverFunction

	//functions that 
	valueReceiverFunctions []ReceiverFunction

	//structs included anonymously in this struct
	inheritedStructs []Struct
}


//for if struct is found as an Anonymous member of something else first
func makeStructUnknown(b *BaseType) {

}


//possibilities for lines:
//Type -> inheritedStruct
//(comma seperated list of names) Type -> NameTypePairs
func makeStruct(b *BaseType, lines []string) *Struct {
	s := Struct{b, nil, nil, nil, nil}
	b.node = &s

	return &s
}