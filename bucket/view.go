package bucket

import (
	"github.com/juelko/bucket/pkg/errors"
	"github.com/juelko/bucket/pkg/events"
)

func NewView(id events.EntityID, stream ...events.Event) (*View, error) {
	const op errors.Op = "bucket.NewView"

	s, err := buildState(id, stream)
	if err != nil {
		return nil, errors.New(op, errors.KindUnexpected, "Error when build state for view", err)
	}

	return &View{
		ID:          string(s.id),
		Title:       string(s.title),
		Description: string(s.desc),
		Version:     uint(s.v),
		IsClosed:    s.closed,
	}, nil
}

type View struct {
	ID          string
	Title       string
	Description string
	Version     uint
	IsClosed    bool
}
