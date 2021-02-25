package bucket

import (
	stderr "errors"
	"time"

	"github.com/juelko/bucket/pkg/errors"
)

var (
	ErrBucketIsClosed error = stderr.New("Bucket is closed")
)

type Event interface {
	StreamID() string
	Type() string
	OccuredAt() time.Time
	Version() uint
	RequestID() string
}

// BaseRequest holds common field for all requests
type BaseRequest struct {
	ID        string `validate:"alphanum,min=3,max=64"`
	RequestID string `validate:"uuid4_rfc4122"`
}

// BaseEvent has common fields and methods to all domain events.
// When BaseEvent is embedded, Domain Event only needs Type() method to satisfy Event interface
type BaseEvent struct {
	ID      string    `validate:"alphanum,min=3,max=64"`
	Occured time.Time `validate:"required" deep:"-"`
	ReqID   string    `validate:"uuid4_rfc4122"`
	V       uint      `validate:"gt=0"`
}

func (be BaseEvent) StreamID() string {
	return be.ID
}

func (be BaseEvent) Version() uint {
	return be.V
}

func (be BaseEvent) RequestID() string {
	return be.ReqID
}

func (be BaseEvent) OccuredAt() time.Time {
	return be.Occured
}

type OpenRequest struct {
	BaseRequest
	Name        string `validate:"alphanum,min=3,max=64"`
	Description string
}

// Opened is a domain event and is emitted when bucket is opened
type Opened struct {
	BaseEvent
	Name        string `validate:"alphanum,min=3,max=64"`
	Description string
}

func (e Opened) Type() string {
	return "bucket.Opened"
}

func Open(req *OpenRequest) Event {
	const op errors.Op = "bucket.Open"

	if err := validateStruct(req); err != nil {
		return newErrorEvent(errors.New(op, errors.KindValidation, "Request validation failed", err), req.RequestID)
	}

	be := BaseEvent{
		ID:      req.ID,
		Occured: time.Now(),
		ReqID:   req.RequestID,
		V:       1,
	}

	return Opened{be, req.Name, req.Description}
}

type CloseRequest struct {
	BaseRequest
}

type Closed struct {
	BaseEvent
}

func (e Closed) Type() string {
	return "bucket.Closed"
}

func Close(req *CloseRequest, stream []Event) Event {
	const op errors.Op = "bucket.Close"

	if err := validateStruct(req); err != nil {
		return newErrorEvent(errors.New(op, errors.KindValidation, err), req.RequestID)
	}

	s, err := newState(req.ID, stream)
	if err != nil {
		return newErrorEvent(errors.New(op, errors.KindUnexpected, err), req.RequestID)
	}

	if s.closed {
		err := errors.New("Closed bucket")
		return newErrorEvent(err, req.RequestID)
	}

	be := BaseEvent{
		ID:      req.ID,
		Occured: time.Now(),
		V:       s.v + 1,
		ReqID:   req.RequestID,
	}

	return Closed{be}

}

type UpdateRequest struct {
	BaseRequest
	Name        string `validate:"alphanum,min=3,max=64"`
	Description string
}

type Updated struct {
	BaseEvent
	Name        string
	Description string
}

func (e Updated) Type() string {
	return "bucket.Updated"
}

func Update(req UpdateRequest, stream []Event) Event {
	const op errors.Op = "bucket.Update"

	if err := validate.Struct(req); err != nil {
		return newErrorEvent(errors.New(op, errors.KindValidation, err), req.RequestID)
	}

	s, err := newState(req.ID, stream)
	if err != nil {
		return newErrorEvent(errors.New(op, errors.KindUnexpected, err), req.RequestID)
	}

	if s.closed {
		err := errors.New("Closed bucket")
		return newErrorEvent(err, req.RequestID)
	}

	be := BaseEvent{
		ID:      req.ID,
		Occured: time.Now(),
		V:       s.v + 1,
		ReqID:   req.RequestID,
	}

	return Updated{be, req.Name, req.Description}
}

type ErrorEvent struct {
	BaseEvent
	Err *errors.Error
}

func (ee ErrorEvent) Type() string {
	return "bucket.ErrorEvent"
}

func (ee ErrorEvent) Error() string {
	return ee.Err.Error()
}

func newErrorEvent(err *errors.Error, rid string) Event {

	be := BaseEvent{
		Occured: time.Now(),
		ReqID:   rid,
	}
	return ErrorEvent{be, err}
}
