package many

type NotImplementing struct {
	id string
}

func (i *NotImplementing) ID() string {
	return i.id
}
