package main

import "github.com/firegoblin/gographviz"

//represents any type without pointer
type BaseType struct {
	//string representation of type
	name string

	//where the BaseType's information is stored
	//may be nil
	node gographviz.GraphableNode

	pkgName string
}

//type handles associating allLevels
func makeBase(s string, pkg string) *BaseType {
	x := BaseType{s, nil, pkg}
	return &x
}

func (b *BaseType) addNode(n gographviz.GraphableNode) {
	b.node = n
}

func (b BaseType) String() string {
	retval := b.pkgName
	if retval == *defaultPackageName {
		retval = ""
	}
	if retval != "" {
		retval += " "
	}
	retval += b.name
	return retval
}
