package main

import "go/ast"

type receiverFunction struct {
	f             *function
	receiverField namedType
}

func newReceiverFunction(f *function, field *ast.Field) receiverFunction {
	return receiverFunction{f, newNamedTypeFromField(field)}
}

func (r *receiverFunction) SlimString() string {
	return r.f.String()
}

func (r *receiverFunction) String() string {
	return "func (" + r.receiverField.String() + ") " + r.f.String()
}
