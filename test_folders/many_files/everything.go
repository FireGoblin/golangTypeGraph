package many

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
