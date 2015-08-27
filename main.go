package main

var builtinTypes = [...]string{"bool", "byte", "complex64", "complex128", "error", "float32", "float64", 
										"int", "int8", "int16", "int32", "int64",
										"rune", "string", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr"}
var typeMap = MasterTypeMap(make(map[string]*Type))
var funcMap = MasterFuncMap(make(map[string]*Function))

func main() {
	//initialize master map with builtin types
	for _, v := range builtinTypes {
		typeMap.lookupOrAdd(v)
	}

	//interfaces found
	//interfaceList := make([]*Interface, 0)
	//structs found := make([]*Struct, 0)
}