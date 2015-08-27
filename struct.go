package main

import . "regexp"
import "strings"

//may be removable
var AnonymousStructParser = MustCompile(`^[^ ]+$`)

//A node type
type Struct struct {
	target *BaseType

	fields []NamedType

	//receiver functions that only work with pointer to this type
	pointerFunctions []*Function

	//functions that 
	valueFunctions []*Function

	//structs included anonymously in this struct
	inheritedStructs []*Struct

	//structs this type is included in anonymously
	includedIn []*Struct
}


//for if struct is found as an Anonymous member of something else first
func makeStructUnknown(b *BaseType, source *Struct) *Struct {
	retval := Struct{b, make([]NamedType, 0), make([]*Function, 0), make([]*Function, 0), make([]*Struct, 0), make([]*Struct, 0)}
	b.node = &retval

	retval.includedIn = append(retval.includedIn, source)

	return &retval
}


//possibilities for lines:
//Type -> inheritedStruct
//(comma seperated list of names) Type -> NamedTypes
//b: the baseType for this struct
//lines: lines from the structs declaration block, preceeding and trailing whitespace removed
func makeStruct(b *BaseType, lines []string) *Struct {
	retval := Struct{b, make([]NamedType, 0), make([]*Function, 0), make([]*Function, 0), make([]*Struct, 0), make([]*Struct, 0)}
	b.node = &retval

	for _, v := range lines {
		ntp := NamedTypeParser.FindStringSubmatch(v)
		if len(ntp) != 0 {
			str := strings.Split(ntp[1], ", ")
			typ := typeMap.lookupOrAdd(ntp[2])
			for _, s := range str {
				retval.fields = append(retval.fields, NamedType{s, typ})
			}
		} else {
			typ := typeMap.lookupOrAdd(v)
			var struc *Struct

			if typ.base.node == nil {
				struc := makeStructUnknown(typ.base, &retval)
				typ.base.node = struc
			} else {
				var ok bool
				struc, ok = typ.base.node.(*Struct)
				if !ok {
					panic("Could not find struct of anonymous member")
				}
			}
			retval.inheritedStructs = append(retval.inheritedStructs, struc)
		}
	}

	return &retval
}