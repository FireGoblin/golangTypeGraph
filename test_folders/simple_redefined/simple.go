package simple

type OneFunction interface {
	ID() int
}

type TwoFunctions interface {
	CombinedName(string) string
	AddInt(int)
}

type ThreeFunctions interface {
	OneFunction
	TwoFunctions
}

type Implements struct {
	id int
}

func (i *Implements) ID() int {
	return i.id
}

type Redefined Implements
