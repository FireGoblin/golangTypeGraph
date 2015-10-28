package redefined

import "fmt"

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
	implements
	name string
}

func (p *partial) combinedName(s string) string {
	return s + p.name
}

func (p *partial) extraFunction(i int) {
	p.id = p.id + i
}

type everything struct {
	partial
}

func (p *everything) addInt(i int) {
	p.id = p.id + i
}

type loseItAll everything

func (l *loseItAll) id() int {
	return 0
}

type direct struct {
	threeFunctions
}

type indirect struct {
	oneFunction
	twoFunctions
}

type mixed struct {
	implements
	twoFunctions
}

type manual struct {
	twoFunctions
	id int
}

func (m *manual) id() int {
	return m.id
}

type onlyOne struct {
	oneFunction
}

type onlyTwo struct {
	twoFunctions
}
