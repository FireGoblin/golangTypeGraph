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
func newBase(s string, pkg string) *BaseType {
	x := BaseType{s, nil, pkg}
	return &x
}

func (b *BaseType) addNode(n gographviz.GraphableNode) {
	b.node = n
}

//Name is safe for use in graph
func (b BaseType) Name() string {
	retval := b.pkgName
	if retval == *defaultPackageName {
		retval = ""
	} else if retval != "" {
		retval += "_SEL_"
	}
	retval += gographviz.SafeName(b.name)
	return retval
}

var builtinTypes = [...]string{"bool", "byte", "complex64", "complex128", "error", "float32", "float64",
	"int", "int8", "int16", "int32", "int64",
	"rune", "string", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr"}

func (b BaseType) StringRelativePkg(pkg string) string {
	retval := b.pkgName
	if retval == pkg {
		retval = ""
	}

	found := false
	for _, v := range builtinTypes {
		if v == retval {
			found = true
			break
		}
	}

	if found {
		retval = ""
	}

	if retval != "" {
		retval += "."
	}
	retval += b.name
	return retval
}

func (b BaseType) String() string {
	return b.StringRelativePkg(*defaultPackageName)
}
