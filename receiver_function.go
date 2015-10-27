package main

import "go/ast"

type ReceiverFunction struct {
	f             *Function
	receiverField NamedType
}

func NewReceiverFunction(f *Function, field *ast.Field) ReceiverFunction {
	return ReceiverFunction{f, NamedTypeFromField(field)}
}

func (r *ReceiverFunction) String() string {
	return "func (" + r.receiverField.String() + ") " + r.receiverField.String()
}
