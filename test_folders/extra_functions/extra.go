package extra

import "fmt"

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

type implements struct {
	id int
}

func (i *implements) String() string {
	return fmt.Sprintf("nonsense %d", i.id)
}

func (i *implements) id() int {
	return i.id
}

type notImplementing struct {
	id string
}

func (i *notImplementing) id() string {
	return i.id
}

type partial struct {
	id   int
	name string
}

func (p *partial) id() int {
	return p.id
}

func (p *partial) combinedName(s string) string {
	return s + p.name
}

func (p *partial) extraFunction(i int) {
	p.id = p.id + i
}

type everything struct {
	id   int
	name string
}

func (p *everything) id() int {
	return p.id
}

func (p *everything) combinedName(s string) string {
	return s + p.name
}

func (p *everything) addInt(i int) {
	p.id = p.id + i
}
