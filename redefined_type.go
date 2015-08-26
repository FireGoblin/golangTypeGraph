package main 

import . "regexp"

//note: does not inherit functions from inheritedFrom
//A node type
type RedefinedType struct {
	target BaseType 

	inheritedFrom Type
	ownFunctions []ReceiverFunction
}

func redefinedTypeParser() *Regexp {
	return MustCompile(`^type (.*?) (.*)$`)
}