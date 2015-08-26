package main 

import . "regexp"

//note: does not inherit functions from inheritedFrom
type RedefinedType struct {
	target Type 

	inheritedFrom Type
	ownFunctions []ReceiverFunction
}

func redefinedTypeParser() *Regexp {
	return MustCompile(`^type (.*?) (.*)$`)
}