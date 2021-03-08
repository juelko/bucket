package bucket

import (
	"github.com/juelko/bucket/pkg/errors"
	"github.com/juelko/bucket/pkg/events"
)

type BucketData struct {
	Title       // bucket title
	Description // bucket description
}

// Opened is a domain event and is emitted when bucket is opened
type Opened struct {
	events.Base // Base event
	BucketData  // Bucket data
}

func (e *Opened) Type() string {
	return "bucket.Opened"
}

func (e *Opened) Data() interface{} {
	return e.BucketData
}

// Business logic for opening
func Open(req *OpenRequest) (*Opened, error) {
	const op errors.Op = "bucket.Open"

	if err := req.Validate(); err != nil {
		return nil, errors.New(op, errors.KindValidation, "invalid request", err)
	}

	return &Opened{
		events.Base{ID: req.ID, V: 1},
		BucketData{Title: req.Title, Description: req.Desc},
	}, nil
}

// Updated is a domain event and is emitted when bucket is updated
type Updated struct {
	events.Base // Base event
	BucketData  // Bucket data
}

func (e *Updated) Type() string {
	return "bucket.Updated"
}

func (e *Updated) Data() interface{} {
	return e.BucketData
}

// Business logic for updating
func Update(req *UpdateRequest, stream []events.Event) (events.Event, error) {
	const op errors.Op = "bucket.Update"

	if err := req.Validate(); err != nil {
		return nil, errors.New(op, errors.KindValidation, "invalid request", err)
	}

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

	return &Updated{
		events.Base{ID: req.ID, V: s.v + 1},
		BucketData{Title: req.Title, Description: req.Desc},
	}, nil
}

// Closed is a domain event and is emitted when bucket is closed
type Closed struct {
	events.Base // Base event
}

func (e *Closed) Type() string {
	return "bucket.Closed"
}
func (e *Closed) Data() interface{} {
	return nil
}

// Business logic for closing
func Close(req *CloseRequest, stream []events.Event) (events.Event, error) {
	const op errors.Op = "bucket.Close"

	if err := req.Validate(); err != nil {
		return nil, errors.New(op, errors.KindValidation, "invalid request", err)
	}

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

	return &Closed{
		events.Base{ID: req.ID, V: s.v + 1},
	}, nil
}
