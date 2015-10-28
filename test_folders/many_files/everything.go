package many

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
