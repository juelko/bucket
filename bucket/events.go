package bucket

import (
	"errors"
	"time"
)

type Event interface {
	ID() string
	Kind() string
	OccuredAt() time.Time
	Version() uint
	RequestID() string
}

type OpenRequest struct {
	ID          BucketID
	Name        string
	Description string
	RequestID   string
}

type Opened struct {
	baseEvent
	Name        string
	Description string
}

func (e Opened) Kind() string {
	return "bucket.Opened"
}

func Open(req OpenRequest) Event {

	if err := req.ID.Validate(); err != nil {
		return newErrorEvent(err, req.RequestID)
	}

	if req.Name == "" {
		err := errors.New("Name cannot be empty")
		return newErrorEvent(err, req.RequestID)
	}

	be := baseEvent{
		id:      req.ID,
		occured: time.Now(),
		rid:     req.RequestID,
		version: 1,
	}

	return Opened{be, req.Name, req.Description}
}

type CloseRequest struct {
	ID        BucketID
	RequestID string
}

type Closed struct {
	baseEvent
}

func (e Closed) Kind() string {
	return "bucket.Closed"
}

func Close(req CloseRequest, stream []Event) Event {

	if err := req.ID.Validate(); err != nil {
		return newErrorEvent(err, req.RequestID)
	}

	s := newState(stream)

	if s.version == 0 {
		err := errors.New("Empty stream")
		return newErrorEvent(err, req.RequestID)
	}

	if req.ID == s.id {
		err := errors.New("ID Mismatch")
		return newErrorEvent(err, req.RequestID)
	}

	if s.closed {
		err := errors.New("Closed bucket")
		return newErrorEvent(err, req.RequestID)
	}

	be := baseEvent{
		id:      req.ID,
		occured: time.Now(),
		version: s.version + 1,
		rid:     req.RequestID,
	}

	return Closed{be}

}

type UpdateRequest struct {
	ID          BucketID
	Name        string
	Description string
	RequestID   string
}

func Update(req UpdateRequest, stream []Event) Event {

	if err := req.ID.Validate(); err != nil {
		return newErrorEvent(err, req.RequestID)
	}

	if req.Name == "" {
		err := errors.New("Name cannot be empty")
		return newErrorEvent(err, req.RequestID)
	}

	s := newState(stream)

	if s.version == 0 {
		err := errors.New("Empty stream")
		return newErrorEvent(err, req.RequestID)
	}

	if req.ID == s.id {
		err := errors.New("ID Mismatch")
		return newErrorEvent(err, req.RequestID)
	}

	if s.closed {
		err := errors.New("Closed bucket")
		return newErrorEvent(err, req.RequestID)
	}

	be := baseEvent{
		id:      req.ID,
		occured: time.Now(),
		version: s.version + 1,
		rid:     req.RequestID,
	}

	return Updated{be, req.Name, req.Description}
}

type Updated struct {
	baseEvent
	Name        string
	Description string
}

func (e Updated) Kind() string {
	return "bucket.Updated"
}

func newErrorEvent(err error, rid string) Event {
	be := baseEvent{
		occured: time.Now(),
		rid:     rid,
	}
	return ErrorEvent{be, err}
}

type ErrorEvent struct {
	baseEvent
	err error
}

func (ee ErrorEvent) Kind() string {
	return "bucket.ErrorEvent"
}

func (ee ErrorEvent) Error() error {
	return ee.err
}

type baseEvent struct {
	id      BucketID
	occured time.Time
	version uint
	rid     string
}

func (be baseEvent) ID() string {
	return be.id.String()
}

func (be baseEvent) OccuredAt() time.Time {
	return be.occured
}

func (be baseEvent) Version() uint {
	return be.version
}

func (be baseEvent) RequestID() string {
	return be.rid
}
