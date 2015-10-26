package main

import "github.com/firegoblin/gographviz"

//data type for fields/parameters.
//A pair of the name of the variable and its type.
type NamedType struct {
	name   string
	target *Type
}

func (n NamedType) String() string {
	return n.name + "\t" + n.target.String()
}

func (n NamedType) Node() gographviz.GraphableNode {
	return n.target.base.node
}
