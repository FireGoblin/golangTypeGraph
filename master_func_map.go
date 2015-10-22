package main

//master map uses singleton pattern, only one of them should be created in the program
//only master should call creators for types
type MasterFuncMap map[string]*Function

func (m MasterFuncMap) lookupOrAdd(s string) *Function {
	x, ok := m[s]

	if ok {
		return x
	} else {
		//m[s] = makeType(s)
		err := makeFunction(s)
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
