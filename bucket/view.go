package bucket

import (
	"time"

	"github.com/juelko/bucket/pkg/errors"
)

func NewView(id ID, stream []Event) (View, error) {
	const op errors.Op = "bucket.NewView"

	s, err := buildState(id, stream)
	if err != nil {
		return View{}, errors.New(op, errors.KindUnexpected, "Error when build state for view", err)
	}

	return View{
		ID:          string(s.id),
		Name:        string(s.name),
		Description: string(s.desc),
		Version:     uint(s.v),
		IsClosed:    s.closed,
		LastUpdate:  s.last,
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
