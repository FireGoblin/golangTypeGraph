package many

type OneFunction interface {
	ID() int
}

type TwoFunctions interface {
	CombinedName(string) string
	AddInt(int)
}

func AddXToY(x OneFunction, y TwoFunctions) {
	y.AddInt(x.ID())
}
