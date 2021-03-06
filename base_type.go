package main

import "github.com/firegoblin/gographviz"

//represents any type without pointer
type baseType struct {
	//string representation of type
	name string

	//where the baseType's information is stored
	//may be nil
	node gographviz.GraphableNode

	pkgName string
}

//type handles associating allLevels
func newBase(s string, pkg string) *baseType {
	x := baseType{s, nil, pkg}
	return &x
}

func (b *baseType) addNode(n gographviz.GraphableNode) {
	b.node = n
}

//used to replace . which doesn't work for a dot graph node name
const periodReplacement = "_SEL_"

//Name is safe for use in graph
func (b baseType) Name() string {
	retval := b.pkgName
	if retval == *defaultPackageName {
		retval = ""
	} else if retval != "" {
		retval += periodReplacement
	}
	retval += gographviz.SafeName(b.name)
	return retval
}

var builtinTypes = [...]string{"bool", "byte", "complex64", "complex128", "error", "float32", "float64",
	"int", "int8", "int16", "int32", "int64",
	"rune", "string", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr"}

func isBuiltinType(s string) bool {
	for _, v := range builtinTypes {
		if v == s {
			return true
		}
	}
	return false
}

func (b baseType) stringRelativePkg(pkg string) string {
	retval := b.pkgName
	if retval == pkg {
		retval = ""
	}

	if isBuiltinType(retval) {
		retval = ""
	}

	if retval != "" {
		retval += "."
	}
	retval += b.name
	return retval
}

func (b baseType) String() string {
	return b.stringRelativePkg(*defaultPackageName)
}
