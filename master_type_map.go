package main

import "go/ast"

//master map uses singleton pattern, only one of them should be created in the program
//only master should call creators for types
type MasterTypeMap struct {
	theMap     map[string]map[string]*Type //map from package name -> type name -> *Type
	currentPkg string                      //used as first index if requester doesn't explicitly pick a package
}

var typeMap = MasterTypeMap{make(map[string]map[string]*Type), ""}

func (m MasterTypeMap) currentMap() map[string]*Type {
	_, ok := m.theMap[m.currentPkg]
	if !ok {
		m.theMap[m.currentPkg] = make(map[string]*Type)
	}
	return m.theMap[m.currentPkg]
}

func (m MasterTypeMap) getPkg(pkg string) map[string]*Type {
	_, ok := m.theMap[pkg]
	if !ok {
		m.theMap[pkg] = make(map[string]*Type)
	}
	return m.theMap[pkg]
}

//TODO: Consider whether to attach the object
//only call if the string is known to be a base type
func (m MasterTypeMap) lookupOrAdd(s string) *Type {
	x, ok := m.currentMap()[s]

	if !ok {
		m.currentMap()[s] = newType(s, m.currentPkg)

		x, ok = m.currentMap()[s]
		if !ok {
			panic("masterlist not properly associated with new type")
		}
	}

	return x
}

func (m MasterTypeMap) lookupOrAddWithPkg(s string, pkg string) *Type {
	x, ok := m.getPkg(pkg)[s]

	if !ok {
		m.getPkg(pkg)[s] = newType(s, pkg)

		x, ok = m.getPkg(pkg)[s]
		if !ok {
			panic("masterlist not properly associated with new type")
		}
	}

	return x
}

func (m MasterTypeMap) lookupOrAddFromExpr(expr ast.Expr) *Type {
	selector, pkg := ReplaceSelector(expr)
	targetPkg := m.currentPkg
	targetExpr := expr
	if pkg != nil {
		targetPkg = String(pkg)
		targetExpr = selector
	}

	s := String(targetExpr)

	x, ok := m.getPkg(targetPkg)[s]

	if ok && x.astNode == nil {
		x.astNode = targetExpr
	} else if !ok {
		m.getPkg(targetPkg)[s] = newTypeFromExpr(targetExpr, targetPkg)

		//error checking
		x, ok = m.getPkg(targetPkg)[s]
		if !ok || x == nil {
			panic("masterlist not properly associated with new func")
		}
	}

	return x
}
