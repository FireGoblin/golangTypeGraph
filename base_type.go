package main

import "github.com/firegoblin/gographviz"

const maxPointerLevel int = 4

//represents any type without pointer
type BaseType struct {
	//string representation of type
	name string

	//where the BaseType's information is stored
	//may be nil (if it is a function type)
	node gographviz.GraphableNode

	//allLevels should be in order
	//i.e: index in slice = pointerLevel of type
	//can go up to 5 references, ****T
	allLevels [maxPointerLevel + 1]*Type
}

//type handles associating allLevels
func makeBase(s string) *BaseType {
	x := BaseType{s, nil, [maxPointerLevel + 1]*Type{}}
	return &x
}

func (b *BaseType) addNode(n gographviz.GraphableNode) {
	b.node = n
}

func (b *BaseType) addType(t *Type) {
	b.allLevels[t.pointerLevel] = t
}

func (b BaseType) String() string {
	return b.name
}

//UNUSED
func (b BaseType) maxReference() int {
	for i, v := range b.allLevels {
		if v == nil {
			return i - 1
		}
	}

	return maxPointerLevel
}
