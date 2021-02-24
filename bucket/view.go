package bucket

import "time"

func NewView(stream []Event) View {
	s := newState(stream)

	return View{
		ID:          s.id.String(),
		Name:        s.name,
		Description: s.description,
		Version:     s.version,
		IsClosed:    s.closed,
		LastUpdate:  s.updated,
	}
}

type View struct {
	ID          string
	Name        string
	Description string
	Version     uint
	IsClosed    bool
	LastUpdate  time.Time
}
