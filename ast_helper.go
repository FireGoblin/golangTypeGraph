package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strings"
)

//ordering is important
func Normalized(f *ast.FieldList) *ast.FieldList {
	return RemoveNames(Flattened(f))
}

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

func BaseTypeOf(expr ast.Expr) ast.Expr {
	switch e := expr.(type) {
	case *ast.StarExpr:
		return BaseTypeOf(e.X)
	case *ast.ChanType:
		return BaseTypeOf(e.Value)
	case *ast.MapType:
		//TODO: Better system for breaking up map types
		return BaseTypeOf(e.Value)
	case *ast.ArrayType:
		return BaseTypeOf(e.Elt)
	}

	return expr
}

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

func stringWithPkg(pkg string, expr ast.Expr) string {
	return String(insertPkg(pkg, expr))
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
		//use StringInterfaceField for interface fields (x(int)string)
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

//*ast.Field in String() assumes a parameter style Field, this works for interface field
func StringInterfaceField(name string, expr *ast.FuncType) string {
	return name + String(expr.Params) + " " + String(expr.Results)
}
