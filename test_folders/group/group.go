package shuffled

import "fmt"

func AddXToY(x OneFunction, y TwoFunctions) {
	y.AddInt(x.ID())
}

func (i *Implements) String() string {
	return fmt.Sprintf("nonsense %d", i.id)
}

func (i *Implements) ID() int {
	return i.id
}

func (i *NotImplementing) ID() string {
	return i.id
}

func (p *Partial) CombinedName(s string) string {
	return s + p.name
}

func (p *Partial) ExtraFunction(i int) {
	p.id = p.id + i
}

func (p *Everything) AddInt(i int) {
	p.id = p.id + i
}

func (l *LoseItAll) ID() int {
	return 0
}

type (
	Everything struct {
		Partial
	}
	NotImplementing struct {
		id string
	}
	ThreeFunctions interface {
		OneFunction
		TwoFunctions
	}
	Implements struct {
		id int
	}
	OneFunction interface {
		ID() int
	}
	LoseItAll    Everything
	TwoFunctions interface {
		CombinedName(string) string
		AddInt(int)
	}
	Partial struct {
		Implements
		name string
	}
)
