package main

import "go/ast"
import "go/token"
import "strings"
import "reflect"
import "fmt"

//ordering is important
func normalized(f *ast.FieldList) *ast.FieldList {
	return removeNames(flattened(f))
}

func flattened(f *ast.FieldList) *ast.FieldList {
	if f == nil {
		return nil
	}

	x := &ast.FieldList{f.Opening, make([]*ast.Field, 0, len(f.List)*2), f.Closing}

	for _, field := range f.List {
		if len(field.Names) > 1 {
			for j := 0; j < len(field.Names); j++ {
				x.List = append(x.List, field)
				x.List[len(x.List)-1].Names = []*ast.Ident{field.Names[j]}
			}
		} else {
			x.List = append(x.List, field)
		}
	}

	return x
}

func removeNames(f *ast.FieldList) *ast.FieldList {
	if f == nil {
		return nil
	}

	x := &ast.FieldList{f.Opening, make([]*ast.Field, len(f.List)), f.Closing}
	copy(x.List, f.List)
	for i := range x.List {
		x.List[i].Names = nil
	}

	return x
}

func FuncField(f *ast.FuncDecl) *ast.Field {
	return &ast.Field{nil, []*ast.Ident{f.Name}, f.Type, nil, nil}
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
			return "chan<-" + String(e.Value)
		case ast.RECV:
			return "<-chan" + String(e.Value)
		default:
			return "chan" + String(e.Value)
		}
	case *ast.FuncType:
		return "func" + String(e.Params) + " " + String(e.Results)
	case *ast.InterfaceType:
		x := ""
		for _, v := range e.Methods.List {
			x += String(v) + "\n"
		}
		return x
	case *ast.MapType:
		return "map[" + String(e.Key) + "]" + String(e.Value)
	case *ast.StructType:
		x := ""
		for _, v := range e.Fields.List {
			x += String(v) + "\n"
		}
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
		x += " "
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
