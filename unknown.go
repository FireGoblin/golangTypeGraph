package main

import "go/ast"
import "github.com/firegoblin/gographviz"

//import "fmt"

//A node is responsible for the incoming edges to it

//A node type
//implements gographviz.GraphableNode
type Unknown struct {
	target *BaseType

	//structs this type is included in anonymously
	includedIn []gographviz.GraphableNode
}

func (u *Unknown) Name() string {
	return u.target.name
}

func (u *Unknown) Attrs() gographviz.Attrs {
	return nil
}

func (u *Unknown) Edges() []*gographviz.Edge {
	return nil
}

func (u *Unknown) addToIncludedIn(x gographviz.GraphableNode) {
	u.includedIn = append(u.includedIn, x)
}

func (u *Unknown) remakeStruct(spec *ast.TypeSpec) *Struct {
	retval := &Struct{u.target, nil, make([]NamedType, 0), make([]*Function, 0), make([]*BaseType, 0), make([]*Struct, len(u.includedIn)), nil, nil, spec.Type}
	copy(retval.includedIn, []*Struct(u.includedIn))

	switch t := spec.Type.(type) {
	case *ast.StructType:
		//struct
		flattenedFields := flattened(t.Fields)
		for _, v := range flattenedFields.List {
			if len(v.Names) != 0 {
				retval.fields = append(retval.fields, NamedType{v.Names[0].Name, typeMap.lookupOrAdd(String(v.Type))})
			} else {
				lookup := typeMap.lookupOrAdd(String(v.Type))
				if lookup.base.node != nil {
					retval.inheritedTypes = append(retval.inheritedTypes, lookup.base)
					lookup.base.node.addToIncludedIn(u)
				} else {
					retval.inheritedTypes = append(retval.inheritedTypes, makeUnknown(retval, lookup.base).target)
				}
			}
		}
	case *ast.Ident:
		//redefined type
		retval.parent = typeMap.lookupOrAdd(t.Name)
	default:
		panic("unexpected type in makeStruct")
	}

	retval.target.addNode(retval)
	return retval
}

func (u *Unknown) remakeInterface(spec *ast.TypeSpec) *Interface {
	interfaceType, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		panic("bad ast.TypeSpec that is not InterfaceType in makeInterface")
	}

	retval := &Interface{u.target, make([]*Function, 0), make([]*Interface, 0), make([]gographviz.GraphableNode, len(u.includedIn)), nil, nil, interfaceType}
	copy(retval, u.includedIn)

	for _, v := range interfaceType.Methods.List {
		if len(v.Names) != 0 {
			f := funcMap.lookupOrAddFromExpr(v.Names[0].Name, v.Type.(*ast.FuncType))
			f.addInterface(retval)
			retval.requiredFunctions = append(retval.requiredFunctions, f)
		} else {
			lookup := typeMap.lookupOrAdd(String(v.Type))
			node := lookup.base.node
			if node != nil {
				retval.inheritedInterfaces = append(retval.inheritedInterfaces, node.(*Interface))
				node.(*Interface).addToIncludedIn(retval)
			} else {
				retval.inheritedInterfaces = append(retval.inheritedInterfaces, makeInterfaceUnknown(retval, lookup.base))
			}
		}
	}

	retval.target.addNode(retval)
	return retval
}
