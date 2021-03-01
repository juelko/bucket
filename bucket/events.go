package bucket

import (
	"time"

	"github.com/juelko/bucket/pkg/errors"
	"github.com/juelko/bucket/pkg/validator"
)

type Event interface {
	StreamID() ID
	Type() string
	OccuredAt() time.Time
	Version() Version
}

// NewOpenRequest takes arguments for new OpenRequest and validates them
// If arguments are valid, *OpenRequest and nil error is returned
// If validation of any argument fails, nil and error is returned
func NewOpenRequest(id ID, n Name, d Description) (*OpenRequest, error) {
	const op errors.Op = "bucket.NewOpenRequest"

	if err := validator.Validate(id, n); err != nil {
		return nil, errors.New(op, errors.KindValidation, "Invalid arguments", err)
	}

	return &OpenRequest{id, n, d}, nil

}

// OpenRequest is holds validated arguments for Open
type OpenRequest struct {
	// contains filtered or unexported fields
	id   ID
	name Name
	desc Description
}

// Opened is a domain event and is emitted when bucket is opened
type Opened struct {
	// contains filtered or unexported fields
	baseEvent
	name Name
	desc Description
}

func (e *Opened) Type() string {
	return "bucket.Opened"
}

func (e *Opened) Name() Name {
	return e.name
}

func (e *Opened) Description() Description {
	return e.desc
}

// Open creates Event from OpenRequest.
// If business requirements are met, retuned event is Opened event
// In case of any error, ErrorEvent is returned
func Open(req *OpenRequest) Event {
	return opening(req)
}

func opening(req *OpenRequest) Event {
	be := baseEvent{
		id: req.id,
		at: time.Now(),
		v:  1,
	}

	return &Opened{be, req.name, req.desc}
}

// NewCloseRequest takes arguments for new CloseRequest and validates them
// If arguments are valid, *CloseRequest and nil error is returned
// If validation of any argument fails, nil and error is returned
func NewCloseRequest(id ID) (*CloseRequest, error) {
	const op errors.Op = "bucket.NewCloseRequest"

	if err := validator.Validate(id); err != nil {
		return nil, errors.New(op, errors.KindValidation, "Invalid arguments", err)
	}

	return &CloseRequest{id}, nil

}

type CloseRequest struct {
	// contains filtered or unexported fields
	id ID
}

type Closed struct {
	// contains filtered or unexported fields
	baseEvent
}

func (e *Closed) Type() string {
	return "bucket.Closed"
}

func Close(req *CloseRequest, stream []Event) Event {
	return stateForClosing(req, stream)
}

func stateForClosing(req *CloseRequest, stream []Event) Event {
	const op errors.Op = "bucket.stateForClosing"

	s, err := buildState(req.id, stream)
	if err != nil {
		return NewErrorEvent(errors.New(op, errors.KindUnexpected, "Error when building state for closing", err))
	}

	return closing(req, &s)
}

func closing(req *CloseRequest, s *state) Event {
	const op errors.Op = "bucket.closing"

	if s.closed {
		return NewErrorEvent(errors.New(op, errors.KindExpected, "Bucket allready closed"))
	}

	be := baseEvent{
		id: req.id,
		at: time.Now(),
		v:  s.v + 1,
	}

	return &Closed{be}
}

func NewUpdateRequest(id ID, n Name, d Description) (*UpdateRequest, error) {
	const op errors.Op = "bucket.NewUpdateRequest"

	if err := validator.Validate(id, n); err != nil {
		return nil, errors.New(op, errors.KindValidation, "Invalid arguments", err)
	}

	return &UpdateRequest{id, n, d}, nil

}

type UpdateRequest struct {
	// contains filtered or unexported fields
	id   ID
	name Name
	desc Description
}

type Updated struct {
	// contains filtered or unexported fields
	baseEvent
	name Name
	desc Description
}

func (e *Updated) Type() string {
	return "bucket.Updated"
}

func (e *Updated) Name() Name {
	return e.name
}

func (e *Updated) Description() Description {
	return e.desc
}

func Update(req *UpdateRequest, stream []Event) Event {
	return stateForUpdating(req, stream)
}

func stateForUpdating(req *UpdateRequest, stream []Event) Event {
	const op errors.Op = "bucket.stateForUpdating"

	s, err := buildState(req.id, stream)
	if err != nil {
		return NewErrorEvent(errors.New(op, errors.KindUnexpected, "Error when building state for update", err))
	}

	return updating(req, &s)
}

func updating(req *UpdateRequest, s *state) Event {
	const op errors.Op = "bucket.updating"

	if s.closed {
		return NewErrorEvent(errors.New(op, errors.KindExpected, "Bucket is closed"))
	}

	be := baseEvent{
		id: req.id,
		at: time.Now(),
		v:  s.v + 1,
	}

	return &Updated{be, req.name, req.desc}
}

type ErrorEvent struct {
	// contains filtered or unexported fields
	baseEvent
	err *errors.Error
}

func (ee *ErrorEvent) Type() string {
	return "bucket.ErrorEvent"
}

func (ee *ErrorEvent) Error() string {
	return ee.err.Error()
}

func NewErrorEvent(err *errors.Error) Event {

	be := baseEvent{
		at: time.Now(),
	}
	return &ErrorEvent{be, err}
}

// baseEvent has common fields and methods to all domain events.
// When baseEvent is embedded, Domain Event only needs Type() method to satisfy Event interface
type baseEvent struct {
	id ID
	at time.Time
	v  Version
}

func (be baseEvent) StreamID() ID {
	return be.id
}

func (be baseEvent) Version() Version {
	return be.v
}

func (be baseEvent) OccuredAt() time.Time {
	return be.at
}
