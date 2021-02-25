package bucket

import "time"

func NewView(id string, stream []Event) (View, error) {
	s, err := newState(id, stream)
	if err != nil {
		return View{}, err
	}

	return View{
		ID:          s.id,
		Name:        s.name,
		Description: s.description,
		Version:     s.v,
		IsClosed:    s.closed,
		LastUpdate:  s.updated,
	}, nil
}

type View struct {
	ID          string
	Name        string
	Description string
	Version     uint
	IsClosed    bool
	LastUpdate  time.Time
}
