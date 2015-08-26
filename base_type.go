package main

const maxPointerLevel int = 5

//represents any type without pointer
type BaseType struct {
	//string representation of type
	name string

	//where the BaseType's information is stored
	//may be nil
	node *interface{}

	//allLevels should be in order
	//i.e: index in slice = pointerLevel of type
	//can go up to 5 references, *****T
	allLevels [maxPointerLevel + 1]*Type
}

//type handles associating allLevels
func makeBase(s string) *BaseType{
	x := BaseType{s, nil, [6]*Type{}}
	return &x
}

func (b *BaseType) addNode(n *interface{}) {
	b.node = n
}

func (b *BaseType) addType(t *Type) {
	b.allLevels[t.pointerLevel] = t
}


func (b BaseType) String() string {
	return b.name
}

func (b BaseType) maxReference() int {
	for i, v := range b.allLevels {
		if v == nil {
			return i - 1
		}
	}

	return maxPointerLevel
}