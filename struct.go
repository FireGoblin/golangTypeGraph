package main

import . "regexp"
import "strings"

//may be removable
var AnonymousStructParser = MustCompile(`^[^ ]+$`)

//A node type
type Struct struct {
	target *BaseType

	//if this is not nil, it is a redefined type
	parent *Type

	//fields should only be empty when a redefined type
	fields []NamedType

	receiverFunctions []*Function

	//structs included anonymously in this struct
	inheritedStructs []*Struct

	//structs this type is included in anonymously
	includedIn []*Struct
}

//no mutation
func (s *Struct) allReceiverFunctions() []*Function {
	retval := make([]*Function, len(s.receiverFunctions))
	c := copy(retval, s.receiverFunctions)
	if c != len(s.receiverFunctions) {
		panic("copy failed in allRequiredFunctions")
	}

	for _, v := range s.inheritedStructs {
		retval = append(retval, v.allReceiverFunctions()...)
	}

	return retval
}

//no mutation
func (s *Struct) implementsInterface(i *Interface) bool {
	required := i.allRequiredFunctions()
	have := s.allReceiverFunctions()

	for _, v := range required {
		found := false
		for _, j := range have {
			if j == v {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

//filter pattern
func (s *Struct) interfacesImplemented(i []*Interface) []*Interface {
	var retval []*Interface

	for _, v := range i {
		if s.implementsInterface(v) {
			retval = append(retval, v)
		}
	}
	
	return retval
}

func (s *Struct) isComposite() bool {
	return len(s.inheritedStructs) > 0
}

func (s *Struct) isRedefinedType() bool {
	return s.parent != nil
}

func (s *Struct) getFields() []NamedType {
	if s.parent == nil {
		return s.fields
	}

	n := s.parent.base.node
	if n == nil {
		return nil
	}

	struc, ok := n.(*Struct)
	if ok {
		return struc.fields
	}

	return nil
}

func makeRedefinedTyp(b *BaseType, s string) *Struct {
	retval := Struct{b, typeMap.lookupOrAdd(s), make([]NamedType, 0), make([]*Function, 0), make([]*Struct, 0), make([]*Struct, 0)}

	return &retval
}

//for if struct is found as an Anonymous member of something else first
func makeStructUnknown(b *BaseType, source *Struct) *Struct {
	retval := Struct{b, nil, make([]NamedType, 0), make([]*Function, 0), make([]*Struct, 0), make([]*Struct, 0)}
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
	retval := Struct{b, nil, make([]NamedType, 0), make([]*Function, 0), make([]*Struct, 0), make([]*Struct, 0)}
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