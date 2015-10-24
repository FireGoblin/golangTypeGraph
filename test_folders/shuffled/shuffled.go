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

type Everything struct {
	Partial
}
type NotImplementing struct {
	id string
}
type ThreeFunctions interface {
	OneFunction
	TwoFunctions
}
type Implements struct {
	id int
}
type OneFunction interface {
	ID() int
}
type LoseItAll Everything
type TwoFunctions interface {
	CombinedName(string) string
	AddInt(int)
}
type Partial struct {
	Implements
	name string
}
