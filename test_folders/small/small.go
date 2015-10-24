package small

import "fmt"

type OneFunction interface {
	ID() int
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
