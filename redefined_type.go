package main 

import . "regexp"

var RedefinedTypeParser = MustCompile(`^type (.*?) (.*)$`)

//note: does not inherit functions from inheritedFrom
//A node type
type RedefinedType struct {
	target BaseType 

	inheritedFrom Type
	ownFunctions []ReceiverFunction
}