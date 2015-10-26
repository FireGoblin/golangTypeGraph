package main

import "github.com/firegoblin/gographviz"

//represents any type without pointer
type BaseType struct {
	//string representation of type
	name string

	//where the BaseType's information is stored
	//may be nil
	node gographviz.GraphableNode
}

//type handles associating allLevels
func makeBase(s string) *BaseType {
	x := BaseType{s, nil}
	return &x
}

func (b *BaseType) addNode(n gographviz.GraphableNode) {
	b.node = n
}

func (b BaseType) String() string {
	return b.name
}
