package main

import "go/ast"
import "go/token"

func removeNames(f *ast.FieldList) *ast.FieldList {
	x := &ast.FieldList{f.Opening, make([]*ast.Field, 0, len(f.List)*2), f.Closing}
	copy(x.List, f.List)
	offset := 0
	for i, field := range f.List {
		if len(field.Names) >= 1 {
			extras := len(field.Names) - 1
			x.List[i+offset].Names = x.List[i+offset].Names[:0]
			for j := 1; j <= extras; j++ {
				x.List = append(x.List[:i+offset+j], append([]*ast.Field{x.List[i+offset]}, x.List[i+offset+j])...)
			}
			offset += extras
		}
	}

	return x
}

// func (f *ast.FuncType) basicFuncType() *ast.FuncType {

// }

func FuncField(f *ast.FuncDecl) *ast.Field {
	return &ast.Field{nil, []*ast.Ident{f.Name}, f.Type, nil, nil}
}

func String(expr ast.Node) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return "*" + String(e.X)
	case *ast.ArrayType:
		switch sub := e.Len.(type) {
		case *ast.BasicLit:
			return "[" + sub.Value + "]" + String(e.Elt)
		case nil:
			return "[]" + String(e.Elt)
		default:
			panic("unexpected type of ast.Expr in ast.ArrayType.Len")
		}
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
		return String(e.Params) + " " + String(e.Results)
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
		x := ""
		if e.Opening != token.NoPos {
			x += ")"
		}
		for _, v := range e.List {
			x += String(v) + ", "
		}
		//TODO: remove last comma
		if e.Closing != token.NoPos {
			x += ")"
		}
		return x
	case *ast.Field:
		x := ""
		for _, v := range e.Names {
			x += v.Name + ", "
		}
		//TODO: remove last comma
		// if reflect.TypeOf(e.Type) != *ast.FuncType {
		// 	x += "/t"
		// }
		x += " "
		x += String(e.Type)
		return x
	default:
		panic("unexpected type of ast.Expr called String()")
	}
}
