package many

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
