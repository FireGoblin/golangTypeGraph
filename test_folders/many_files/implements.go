package many

import "fmt"

type implements struct {
	id int
}

func (i *implements) String() string {
	return fmt.Sprintf("nonsense %d", i.id)
}

func (i *implements) id() int {
	return i.id
}
