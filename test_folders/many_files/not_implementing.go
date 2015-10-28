package many

type notImplementing struct {
	id string
}

func (i *notImplementing) id() string {
	return i.id
}
