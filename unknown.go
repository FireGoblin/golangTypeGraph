package main

import "go/ast"
import "github.com/firegoblin/gographviz"

//A node is responsible for the incoming edges to it

//A node type
//implements gographviz.GraphableNode
type Unknown struct {
	target *BaseType
}

func (u *Unknown) Name() string {
	return gographviz.SafeName(u.target.name)
}

func (u *Unknown) Attrs() gographviz.Attrs {
	return nil
}

func (u *Unknown) Edges() []*gographviz.Edge {
	return nil
}

func (u *Unknown) remakeStruct(spec *ast.TypeSpec) *Struct {
	retval := &Struct{u.target, nil, make([]NamedType, 0), make([]*Function, 0), make([]*BaseType, 0), nil, nil, spec.Type}

	switch t := spec.Type.(type) {
	case *ast.StructType:
		//struct
		flattenedFields := flattened(t.Fields)
		for _, v := range flattenedFields.List {
			if len(v.Names) != 0 {
				retval.fields = append(retval.fields, NamedType{v.Names[0].Name, typeMap.lookupOrAddFromExpr(v.Type)})
			} else {
				lookup := typeMap.lookupOrAddFromExpr(v.Type)
				if lookup.base.node != nil {
					retval.inheritedTypes = append(retval.inheritedTypes, lookup.base)
				} else {
					retval.inheritedTypes = append(retval.inheritedTypes, makeUnknown(retval, lookup.base).target)
				}
			}
		}
	default:
		//redefined type
		retval.parent = typeMap.lookupOrAddFromExpr(t)
	}

	retval.target.addNode(retval)
	return retval
}

func (u *Unknown) remakeInterface(spec *ast.TypeSpec) *Interface {
	interfaceType, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		panic("bad ast.TypeSpec that is not InterfaceType in makeInterface")
	}

	retval := &Interface{u.target, make([]*Function, 0), make([]*Interface, 0), nil, nil, interfaceType}

	for _, v := range interfaceType.Methods.List {
		if len(v.Names) != 0 {
			f := funcMap.lookupOrAddFromExpr(v.Names[0].Name, v.Type.(*ast.FuncType))
			f.addInterface(retval)
			retval.requiredFunctions = append(retval.requiredFunctions, f)
		} else {
			lookup := typeMap.lookupOrAddFromExpr(v.Type)
			node := lookup.base.node
			if node != nil {
				retval.inheritedInterfaces = append(retval.inheritedInterfaces, node.(*Interface))
			} else {
				retval.inheritedInterfaces = append(retval.inheritedInterfaces, makeInterfaceUnknown(retval, lookup.base))
			}
		}
	}

	retval.target.addNode(retval)
	return retval
}
