package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strings"
)

// Normalized is a function to return the same FieldList for any function
// with a matching signature.  exp: Normalized(func(x,y int)) == Normalized(func(a int, b int))
func Normalized(f *ast.FieldList) *ast.FieldList {
	return RemoveNames(Flattened(f))
}

// Flattened converts any fields with multiple names into multiple fields
func Flattened(f *ast.FieldList) *ast.FieldList {
	if f == nil {
		return nil
	}

	x := &ast.FieldList{f.Opening, make([]*ast.Field, 0, len(f.List)*2), f.Closing}

	for _, field := range f.List {
		if len(field.Names) > 1 {
			for j := 0; j < len(field.Names); j++ {
				local := *field
				local.Names = []*ast.Ident{ast.NewIdent(field.Names[j].Name)}
				x.List = append(x.List, &local)
			}
		} else {
			local := *field
			x.List = append(x.List, &local)
		}
	}

	return x
}

// RemoveNames removes the names from all fields in the FieldList
func RemoveNames(f *ast.FieldList) *ast.FieldList {
	if f == nil {
		return nil
	}

	x := &ast.FieldList{f.Opening, make([]*ast.Field, 0, len(f.List)), f.Closing}

	for _, field := range f.List {
		local := *field
		local.Names = nil
		x.List = append(x.List, &local)
	}

	return x
}

// RootTypeOf recurses through the expr to find the base type of the expression
func RootTypeOf(expr ast.Expr) ast.Expr {
	switch e := expr.(type) {
	case *ast.StarExpr:
		return RootTypeOf(e.X)
	case *ast.ChanType:
		return RootTypeOf(e.Value)
	case *ast.MapType:
		//TODO: Better system for breaking up map types
		return RootTypeOf(e.Value)
	case *ast.ArrayType:
		return RootTypeOf(e.Elt)
	}

	return expr
}

// RecursiveTypeOf removes the top level of the expr and returns the result
func RecursiveTypeOf(expr ast.Expr) ast.Expr {
	switch e := expr.(type) {
	case *ast.StarExpr:
		return e.X
	case *ast.ChanType:
		return e.Value
	case *ast.MapType:
		//TODO: Better system for breaking up map types
		return e.Value
	case *ast.ArrayType:
		return e.Elt
	}

	return nil
}

// ReplaceSelector replaces the SelectorExpr in expr with its Sel field
func ReplaceSelector(expr ast.Expr) (replaced ast.Expr, X *ast.Ident) {
	switch e := expr.(type) {
	case *ast.StarExpr:
		local := *e
		local.X, X = ReplaceSelector(e.X)
		return &local, X
	case *ast.ChanType:
		local := *e
		local.Value, X = ReplaceSelector(e.Value)
		return &local, X
	case *ast.MapType:
		local := *e
		local.Value, X = ReplaceSelector(e.Value)
		return &local, X
	case *ast.ArrayType:
		local := *e
		local.Elt, X = ReplaceSelector(e.Elt)
		return &local, X
	case *ast.SelectorExpr:
		return e.Sel, e.X.(*ast.Ident)
	}

	return expr, nil
}

func insertPkg(pkg string, expr ast.Expr) ast.Expr {
	switch e := expr.(type) {
	case *ast.StarExpr:
		local := *e
		local.X = insertPkg(pkg, e.X)
		return &local
	case *ast.ChanType:
		local := *e
		local.Value = insertPkg(pkg, e.Value)
		return &local
	case *ast.MapType:
		local := *e
		local.Value = insertPkg(pkg, e.Value)
		return &local
	case *ast.ArrayType:
		local := *e
		local.Elt = insertPkg(pkg, e.Elt)
		return &local
	case *ast.Ident:
		return &ast.SelectorExpr{ast.NewIdent(pkg), e}
	}

	//returning nil is an error
	return nil
}

func stringWithPkg(pkg string, expr ast.Expr) string {
	return String(insertPkg(pkg, expr))
}

// String uses recursion to print the ast.Node as it would appear in code
func String(expr ast.Node) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return "*" + String(e.X)
	case *ast.SelectorExpr:
		return String(e.X) + "." + String(e.Sel)
	case *ast.BasicLit:
		return e.Value
	case *ast.Ellipsis:
		return "..."
	case *ast.BinaryExpr:
		return String(e.X) + " " + e.Op.String() + " " + String(e.Y)
	case *ast.ArrayType:
		return "[" + String(e.Len) + "]" + String(e.Elt)
	case *ast.ChanType:
		switch e.Dir {
		case ast.SEND:
			return "chan<- " + String(e.Value)
		case ast.RECV:
			return "<-chan " + String(e.Value)
		default:
			return "chan " + String(e.Value)
		}
	case *ast.FuncType:
		return "func" + String(e.Params) + " " + String(e.Results)
	case *ast.InterfaceType:
		x := "interface\\{"
		for _, v := range e.Methods.List {
			x += String(v) + ", "
		}
		x = strings.Trim(x, ", ")
		x += "\\}"
		return x
	case *ast.MapType:
		return "map[" + String(e.Key) + "]" + String(e.Value)
	case *ast.StructType:
		x := "struct\\{"
		for _, v := range e.Fields.List {
			x += String(v) + ", "
		}
		x = strings.Trim(x, ", ")
		x += "\\}"
		return x
	case *ast.FieldList:
		if e == nil {
			return ""
		}
		x := ""
		if e.Opening != token.NoPos {
			x += "("
		}
		for _, v := range e.List {
			x += String(v) + ", "
		}

		x = strings.Trim(x, ", ")
		if e.Closing != token.NoPos {
			x += ")"
		}
		return x
	case *ast.Field:
		//assumes a parameter style Field (x func(int)string)
		//use stringInterfaceField for interface fields (x(int)string)
		x := ""
		for _, v := range e.Names {
			x += v.Name + ", "
		}

		x = strings.Trim(x, ", ")
		if x != "" {
			x += " "
		}
		x += String(e.Type)
		return x
	case nil:
		return ""
	default:
		panic(fmt.Sprintln("unexpected type of ast.Expr called String()", reflect.TypeOf(expr)))
	}
}

func stringInterfaceField(name string, expr *ast.FuncType) string {
	return name + String(expr.Params) + " " + String(expr.Results)
}
