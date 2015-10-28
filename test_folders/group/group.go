package shuffled

import "fmt"

func addXToY(x oneFunction, y twoFunctions) {
	y.addInt(x.id())
}

func (i *implements) String() string {
	return fmt.Sprintf("nonsense %d", i.id)
}

func (i *implements) id() int {
	return i.id
}

func (i *notImplementing) id() string {
	return i.id
}

func (p *partial) combinedName(s string) string {
	return s + p.name
}

func (p *partial) extraFunction(i int) {
	p.id = p.id + i
}

func (p *everything) addInt(i int) {
	p.id = p.id + i
}

func (l *loseItAll) id() int {
	return 0
}

type (
	everything struct {
		partial
	}
	notImplementing struct {
		id string
	}
	threeFunctions interface {
		oneFunction
		twoFunctions
	}
	implements struct {
		id int
	}
	oneFunction interface {
		id() int
	}
	loseItAll    everything
	twoFunctions interface {
		combinedName(string) string
		addInt(int)
	}
	partial struct {
		implements
		name string
	}
)
