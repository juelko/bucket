package bucket

import (
	"time"

	"github.com/juelko/bucket/pkg/errors"
)

type Event interface {
	StreamID() ID
	Type() string
	OccuredAt() time.Time
	Version() Version
	RequestID() RequestID
}

// NewOpenRequest takes arguments for new OpenRequest and validates them
// If arguments are valid, *OpenRequest and nil error is returned
// If validation of any argument fails, nil and error is returned
func NewOpenRequest(id ID, r RequestID, n Name, d Description) (*OpenRequest, error) {
	const op errors.Op = "bucket.NewOpenRequest"

	if err := validateArgs(id, r, n); err != nil {
		return nil, errors.New(op, errors.KindValidation, "Invalid arguments", err)
	}

	return &OpenRequest{id, r, n, d}, nil

}

// OpenRequest is holds validated arguments for Open
type OpenRequest struct {
	// contains filtered or unexported fields
	id    ID
	rid   RequestID
	name  Name
	dessc Description
}

// Opened is a domain event and is emitted when bucket is opened
type Opened struct {
	// contains filtered or unexported fields
	baseEvent
	name Name
	desc Description
}

func (e Opened) Type() string {
	return "bucket.Opened"
}

// Open creates Event from OpenRequest.
// If business requirements are met, retuned event is Opened event
// In case of any error, ErrorEvent is returned
func Open(req *OpenRequest) Event {
	return opening(req)
}

func opening(req *OpenRequest) Event {
	be := baseEvent{
		id:  req.id,
		at:  time.Now(),
		rid: req.rid,
		v:   1,
	}

	return Opened{be, req.name, req.dessc}
}

// NewClosedRequest takes arguments for new CloseRequest and validates them
// If arguments are valid, *CloseRequest and nil error is returned
// If validation of any argument fails, nil and error is returned
func NewClosedRequest(id ID, r RequestID) (*CloseRequest, error) {
	const op errors.Op = "bucket.NewClosedRequest"

	if err := validateArgs(id, r); err != nil {
		return nil, errors.New(op, errors.KindValidation, "Invalid arguments", err)
	}

	return &CloseRequest{id, r}, nil

}

type CloseRequest struct {
	// contains filtered or unexported fields
	id  ID
	rid RequestID
}

type Closed struct {
	// contains filtered or unexported fields
	baseEvent
}

func (e Closed) Type() string {
	return "bucket.Closed"
}

func Close(req *CloseRequest, stream []Event) Event {
	return stateForClosing(req, stream)
}

func stateForClosing(req *CloseRequest, stream []Event) Event {
	const op errors.Op = "bucket.stateForClosing"

	s, err := buildState(req.id, stream)
	if err != nil {
		return newErrorEvent(errors.New(op, errors.KindUnexpected, "Error when building state for closing", err), req.rid)
	}

	return closing(req, &s)
}

func closing(req *CloseRequest, s *state) Event {
	const op errors.Op = "bucket.closing"

	if s.closed {
		return newErrorEvent(errors.New(op, errors.KindExpected, "Bucket allready closed"), req.rid)
	}

	be := baseEvent{
		id:  req.id,
		at:  time.Now(),
		v:   s.v + 1,
		rid: req.rid,
	}

	return Closed{be}
}

type UpdateRequest struct {
	// contains filtered or unexported fields
	id   ID
	rid  RequestID
	name Name
	desc Description
}
type Updated struct {
	// contains filtered or unexported fields
	baseEvent
	name Name
	desc Description
}

func (e Updated) Type() string {
	return "bucket.Updated"
}

func Update(req *UpdateRequest, stream []Event) Event {
	return stateForUpdating(req, stream)
}

func stateForUpdating(req *UpdateRequest, stream []Event) Event {
	const op errors.Op = "bucket.stateForUpdating"

	s, err := buildState(req.id, stream)
	if err != nil {
		return newErrorEvent(errors.New(op, errors.KindUnexpected, "Error when building state for update", err), req.rid)
	}

	return updating(req, &s)
}

func updating(req *UpdateRequest, s *state) Event {
	const op errors.Op = "bucket.updating"

	if s.closed {
		return newErrorEvent(errors.New(op, errors.KindExpected, "Bucket is closed"), req.rid)
	}

	be := baseEvent{
		id:  req.id,
		at:  time.Now(),
		v:   s.v + 1,
		rid: req.rid,
	}

	return Updated{be, req.name, req.desc}
}

type ErrorEvent struct {
	// contains filtered or unexported fields
	baseEvent
	err *errors.Error
}

func (ee ErrorEvent) Type() string {
	return "bucket.ErrorEvent"
}

func (ee ErrorEvent) Error() string {
	return ee.err.Error()
}

func newErrorEvent(err *errors.Error, rid RequestID) Event {

	be := baseEvent{
		at:  time.Now(),
		rid: rid,
	}
	return ErrorEvent{be, err}
}

// baseEvent has common fields and methods to all domain events.
// When baseEvent is embedded, Domain Event only needs Type() method to satisfy Event interface
type baseEvent struct {
	id  ID
	at  time.Time
	rid RequestID
	v   Version
}

func (be baseEvent) StreamID() ID {
	return be.id
}

func (be baseEvent) Version() Version {
	return be.v
}

func (be baseEvent) RequestID() RequestID {
	return be.rid
}

func (be baseEvent) OccuredAt() time.Time {
	return be.at
}
