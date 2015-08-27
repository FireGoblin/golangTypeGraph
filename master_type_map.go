package main

//master map uses singleton pattern, only one of them should be created in the program
//only master should call creators for types
type MasterTypeMap map[string]*Type

func (m MasterTypeMap) lookupOrAdd(s string) *Type {
	x, ok := m[s]

	if ok {
		return x
	} else {
		//m[s] = makeType(s)
		err := makeType(s)
		if err != nil {
			panic("error creating type for master list")
		}

		x, ok = m[s]
		if !ok {
			panic("masterlist not properly associated with new type")
		}

		return x
	}
}