package small

import "fmt"

type oneFunction interface {
	id() int
}

type implements struct {
	id int
}

func (i *implements) string() string {
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
