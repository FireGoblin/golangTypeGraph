package main

import . "regexp"

var NamedTypeParser = MustCompile(`^(.+?)[ ]+((?:func.+)|(?:[^ ]+))$`)

//data type for fields/parameters.
//A pair of the name of the variable and its type.
type NamedType struct {
	name   string
	target *Type
}
