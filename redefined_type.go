package main 

import . "regexp"

var RedefinedTypeParser = MustCompile(`^type (.*?) (.*)$`)

//note: does not inherit functions from inheritedFrom
//A node type
type RedefinedType struct {
	target *BaseType 

	parent *Type

	ownFunctions []*ReceiverFunction
}

func makeRedefinedTyp(b *BaseType, s string) *RedefinedType {
	retval := RedefinedType{b, typeMap.lookupOrAdd(s), make([]*ReceiverFunction, 0)}

	return &retval
}

func (r RedefinedType) fields() []NamedType {
	n := r.parent.base.node
	if n == nil {
		return nil
	}

	struc, ok := n.(*Struct)
	if ok {
		return struc.fields
	}

	return nil
}