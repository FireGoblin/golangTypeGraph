package simple

type oneFunction interface {
	id() int
}

type twoFunctions interface {
	combinedName(string) string
	addInt(int)
}

type threeFunctions interface {
	oneFunction
	twoFunctions
}

type implements struct {
	id int
}

func (i *implements) id() int {
	return i.id
}

type redefined implements
