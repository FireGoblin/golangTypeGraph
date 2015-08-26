package main

import . "regexp"

var NamedTypeParser = MustCompile(`^(.+?)[ ]+((?:func.+)|(?:[^ ]+))$`)

type NamedType struct {
	name string
	target *Type
}