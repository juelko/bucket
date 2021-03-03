package bucket

import (
	"time"

	"github.com/juelko/bucket/pkg/errors"
	"github.com/juelko/bucket/pkg/events"
)

// Opened is a domain event and is emitted when bucket is opened
type Opened struct {
	events.Base             // Base event
	Title       Title       // bucket title
	Desc        Description // bucket description
}

func (e *Opened) Type() string {
	return "bucket.Opened"
}

// Business logic for opening
func open(req *OpenRequest) events.Event {
	return newOpened(req)
}

func newOpened(req *OpenRequest) events.Event {
	be := events.Base{
		ID:  req.ID,
		RID: req.RID,
		At:  time.Now(),
		V:   1,
	}

	return &Opened{be, req.Title, req.Desc}
}

// Updated is a domain event and is emitted when bucket is updated
type Updated struct {
	events.Base             // Base event
	Title       Title       // new bucket title
	Desc        Description // new bucket description
}

func (e *Updated) Type() string {
	return "bucket.Updated"
}

// Business logic for updating
func update(req *UpdateRequest, stream []events.Event) (events.Event, error) {
	return stateForUpdating(req, stream)
}

func stateForUpdating(req *UpdateRequest, stream []events.Event) (events.Event, error) {
	const op errors.Op = "bucket.stateForUpdating"

	s, err := buildState(req.ID, stream)
	if err != nil {
		return nil, errors.New(op, errors.KindUnexpected, "Error when building state for update", err)
	}

	return newUpdate(req, &s)
}

func newUpdate(req *UpdateRequest, s *state) (events.Event, error) {
	const op errors.Op = "bucket.updating"

	if s.closed {
		return nil, errors.New(op, errors.KindExpected, "Bucket is closed")
	}

	be := events.Base{
		ID:  req.ID,
		RID: req.RID,
		At:  time.Now(),
		V:   s.v + 1,
	}

	return &Updated{be, req.Title, req.Desc}, nil
}

// Closed is a domain event and is emitted when bucket is closed
type Closed struct {
	events.Base // Base event
}

func (e *Closed) Type() string {
	return "bucket.Closed"
}

// Business logic for closing
func close(req *CloseRequest, stream []events.Event) (events.Event, error) {
	return stateForClosing(req, stream)
}

func stateForClosing(req *CloseRequest, stream []events.Event) (events.Event, error) {
	const op errors.Op = "bucket.stateForClosing"

	s, err := buildState(req.ID, stream)
	if err != nil {
		return nil, errors.New(op, errors.KindUnexpected, "Error when building state for closing", err)
	}

	return newClosed(req, &s)
}

func newClosed(req *CloseRequest, s *state) (events.Event, error) {
	const op errors.Op = "bucket.closing"

	if s.closed {
		return nil, errors.New(op, errors.KindExpected, "Bucket allready closed")
	}

	be := events.Base{
		ID:  req.ID,
		RID: req.RID,
		At:  time.Now(),
		V:   s.v + 1,
	}

	return &Closed{be}, nil
}
