package main

//import . "regexp"
import "strings"
import "go/ast"

type Type struct {
	name string
	base *BaseType

	//corresponds to the number of asterisks for the type
	//exp: **string would have pointerLevel = 2
	pointerLevel int

	astNode ast.Expr
}

//TODO: maybe more effecient way to do this
//only call from master_type_map
func makeType(s string) *Type {
	return sharedMakeType(s, nil)
}

//TODO: maybe more effecient way to do this
//only call from master_type_map
func makeTypeFromExpr(expr ast.Expr) *Type {
	return sharedMakeType(String(expr), expr)
}

func sharedMakeType(s string, expr ast.Expr) *Type {
	baseType := strings.Trim(s, "*")
	pLevel := len(s) - len(baseType)

	retval := Type{s, nil, pLevel, expr}

	if pLevel == 0 {
		b := makeBase(baseType)
		retval.base = b
		b.addType(&retval)
	} else {
		b, ok := typeMap[baseType]
		if !ok {
			b = typeMap.lookupOrAdd(baseType)
		}

		retval.base = b.base
		b.base.addType(&retval)

		//create lower type if not created yet
		_, ok = typeMap[s[1:]]
		if !ok {
			makeTypeRecursive(s[1:], retval.base, pLevel-1, expr.(*ast.StarExpr).X)
		}
	}

	return &retval
}

//never call outside of makeType
func makeTypeRecursive(s string, b *BaseType, pLevel int, expr ast.Expr) {
	x := Type{s, b, pLevel, expr}
	typeMap[s] = &x
	b.addType(&x)

	_, ok := typeMap[s[1:]]
	if !ok {
		makeTypeRecursive(s[1:], b, pLevel-1, expr.(*ast.StarExpr).X)
	}
}

func (t Type) String() string {
	return t.name
}

func (t Type) BaseName() string {
	return t.base.name
}

// //***************
// func (t Type) isFunc() bool {
// 	return t.String()[0:4] == "func"
// }

// //**************
// //following functions only apply to types that pass isFunc()
// //*************

// func (f Type) params() ([]*Type, error) {
// 	if !f.isFunc() {
// 		return nil, fmt.Errorf("params called on non-function type")
// 	}

// 	results := FuncTypeParser.FindStringSubmatch(f.name)

// 	retval := make([]*Type, 0, len(results))

// 	for _, str := range strings.Split(results[1], ", ") {
// 		retval = append(retval, typeMap.lookupOrAdd(str))
// 	}

// 	return retval, nil
// }

// func (f Type) returnTypes() ([]*Type, error) {
// 	if !f.isFunc() {
// 		return nil, fmt.Errorf("returnTypes called on non-function type")
// 	}

// 	results := FuncTypeParser.FindStringSubmatch(f.name)

// 	retval := make([]*Type, 0, len(results))

// 	for _, str := range strings.Split(results[2], ", ") {
// 		retval = append(retval, typeMap[str])
// 	}

// 	return retval, nil
// }
