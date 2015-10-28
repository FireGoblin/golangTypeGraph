package many

type oneFunction interface {
	id() int
}

type twoFunctions interface {
	combinedName(string) string
	addInt(int)
}

func addXToY(x oneFunction, y twoFunctions) {
	y.addInt(x.id())
}
