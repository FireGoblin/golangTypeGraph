package main 

import . "regexp"
import "strings"
import "fmt"

//----------------------------------------------

type Type struct {
	name string
	base *BaseType

	//corresponds to the number of asterisks for the type
	//exp: **string would have pointerLevel = 2
	pointerLevel int
}

//TODO: maybe more effecient way to do this
func makeType(s string) error {
	_, ok := typeMap[s]
	if ok {
		return fmt.Errorf("attempt to create already created type")
	}

	baseType := strings.TrimLeft(s, "*")
	pLevel := len(s) - len(baseType)
	retval := Type{s, nil, pLevel}

	if pLevel == 0 {
		bt := makeBase(baseType)
		retval.base = bt
		typeMap[s] = &retval
		bt.addType(&retval)
	} else {
		b, ok := typeMap[baseType]
		if !ok {
			b = typeMap.lookupOrAdd(baseType)
		} 

		retval.base = b.base
		typeMap[s] = &retval
		b.base.allLevels[pLevel] = &retval

		//create lower type if not created yet
		_, ok = typeMap[s[1:]]
		if !ok {
			makeTypeRecursive(s[1:], retval.base, pLevel - 1)
		}
	}

	return nil
}

//never call outside of makeType
func makeTypeRecursive(s string, b *BaseType, pLevel int) {
	x := Type{s, b, pLevel}
	typeMap[s] = &x
	b.addType(&x)

	_, ok := typeMap[s[1:]]
	if !ok {
		makeTypeRecursive(s[1:], b, pLevel - 1)
	}
}

func (t Type) String() string {
	return t.name
}

func (t Type) BaseName() string {
	return t.base.name
}

//**************

func FuncTypeParser() *Regexp {
	return MustCompile(`^func\((.*?)\) (.*)$`)
}

func (t Type) isFunc() bool {
	return t.String()[0:4] == "func"
}

func (f Type) params() ([]*Type, error) {
	if !f.isFunc() {
		return nil, fmt.Errorf("params called on non-function type")
	}

	reg := FuncTypeParser()
	results := reg.FindStringSubmatch(f.name)

	retval := make([]*Type, 0, len(results))

	for _, str := range strings.Split(results[1], ", ") {
		retval = append(retval, typeMap.lookupOrAdd(str))
	}

	return retval, nil
}

func (f Type) returnTypes() ([]*Type, error) {
	if !f.isFunc() {
		return nil, fmt.Errorf("returnTypes called on non-function type")
	}

	reg := FuncTypeParser()
	results := reg.FindStringSubmatch(f.name)

	retval := make([]*Type, 0, len(results))

	for _, str := range strings.Split(results[2], ", ") {
		retval = append(retval, typeMap[str])
	}

	return retval, nil
}