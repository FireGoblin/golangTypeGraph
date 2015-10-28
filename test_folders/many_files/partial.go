package many

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
