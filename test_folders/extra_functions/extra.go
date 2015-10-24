package extra

import "fmt"

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

type Implements struct {
	id int
}

func (i *Implements) String() string {
	return fmt.Sprintf("nonsense %d", i.id)
}

func (i *Implements) ID() int {
	return i.id
}

type NotImplementing struct {
	id string
}

func (i *NotImplementing) ID() string {
	return i.id
}

type Partial struct {
	id   int
	name string
}

func (p *Partial) ID() int {
	return p.id
}

func (p *Partial) CombinedName(s string) string {
	return s + p.name
}

func (p *Partial) ExtraFunction(i int) {
	p.id = p.id + i
}

type Everything struct {
	id   int
	name string
}

func (p *Everything) ID() int {
	return p.id
}

func (p *Everything) CombinedName(s string) string {
	return s + p.name
}

func (p *Everything) AddInt(i int) {
	p.id = p.id + i
}
