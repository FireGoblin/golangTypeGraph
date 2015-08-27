package main

//A node type
type Interface struct {
	target *BaseType

	//any functions required to implement the interface
	//does not include functions inherited indirectly through inheritedInterfaces
	requiredFunctions []*Function

	//interfaces this inherits from
	//if not zero, is composite interface
	inheritedInterfaces []*Interface

	//interfaces this is included int
	includedIn []*Interface
}

//for if interface is found as an Anonymous member of something else first
func makeInterfaceUnknown(b *BaseType, source *Interface) *Interface {
	retval := Interface{b, make([]*Function, 0), make([]*Interface, 0), make([]*Interface, 0)}
	b.node = &retval

	retval.includedIn = append(retval.includedIn, source)

	return &retval
}


//possibilities for lines:
//Type -> inheritedStruct
//(comma seperated list of names) Type -> NamedTypes
//b: the baseType for this struct
//lines: lines from the structs declaration block, preceeding and trailing whitespace removed
func makeInterface(b *BaseType, lines []string) *Interface {
	retval := Interface{b, make([]*Function, 0), make([]*Interface, 0), make([]*Interface, 0)}
	b.node = &retval

	for _, v := range lines {
		ifp := FunctionParser.FindStringSubmatch(v)
		if len(ifp) != 0 {
			f := funcMap.lookupOrAdd(ifp[0])
			f.addInterface(&retval)
			retval.requiredFunctions = append(retval.requiredFunctions, f)
		} else {
			typ := typeMap.lookupOrAdd(v)
			var interfac *Interface

			if typ.base.node == nil {
				interfac := makeInterfaceUnknown(typ.base, &retval)
				typ.base.node = interfac
			} else {
				var ok bool
				interfac, ok = typ.base.node.(*Interface)
				if !ok {
					panic("Could not find struct of anonymous member")
				}
			}
			retval.inheritedInterfaces = append(retval.inheritedInterfaces, interfac)
		}
	}

	return &retval
}