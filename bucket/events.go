package bucket

import (
	"encoding/json"
	"time"

	"github.com/juelko/bucket/pkg/errors"
	"github.com/juelko/bucket/pkg/events"
	"github.com/juelko/bucket/pkg/request"
)

// Opened is a domain event and is emitted when bucket is opened
type Opened struct {
	events.Base `json:"-"`  // Base event
	Title       Title       `json:"title"` // bucket title
	Desc        Description `json:"desc"`  // bucket description
}

func (e *Opened) Type() string {
	return "bucket.Opened"
}

func (e *Opened) Payload() []byte {
	d, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	return d
}

// Use only when decoding serialized events from store
func BuildOpened(id events.StreamID, rid request.ID, at time.Time, v events.Version, t Title, d Description) *Opened {
	return &Opened{events.Base{ID: id, RID: rid, At: at, V: v}, t, d}
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
	events.Base `json:"-"`  // Base event
	Title       Title       `json:"title"` // new bucket title
	Desc        Description `json:"desc"`  // new bucket description
}

func (e *Updated) Type() string {
	return "bucket.Updated"
}

func (e *Updated) Payload() []byte {

	d, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	return d
}

// Use only when decoding serialized events from store
func BuildUpdated(id events.StreamID, rid request.ID, at time.Time, v events.Version, t Title, d Description) *Updated {
	return &Updated{events.Base{ID: id, RID: rid, At: at, V: v}, t, d}
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
	events.Base `json:"-"` // Base event
	IsClosed    bool       `json:"isClosed"`
}

func (e *Closed) Type() string {
	return "bucket.Closed"
}

func (e *Closed) Payload() []byte {

	d, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	return d
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

	return &Closed{be, true}, nil
}

// Use only when decoding serialized events from store
func BuildEvent(base events.Base, eventType string, payload []byte) (events.Event, error) {
	const op errors.Op = "bucket.BuildEvent"

	switch eventType {
	case "bucket.Opened":
		return buildOpened(base, payload)
	case "bucket.Updated":
		return buildUpdated(base, payload)
	case "bucket.Closed":
		return buildClosed(base, payload)
	default:
		return nil, errors.New(op, errors.KindUnexpected, "Unkown eventType")
	}

}

// Use only when decoding serialized events from store
func buildOpened(eb events.Base, payload []byte) (*Opened, error) {
	const op errors.Op = "bucket.buildOpened"

	var o Opened

	err := json.Unmarshal(payload, &o)
	if err != nil {
		return nil, errors.New(op, errors.KindUnexpected, "could not unmarshal opened")
	}

	o.Base = eb

	return &o, nil
}

// Use only when decoding serialized events from store
func buildUpdated(eb events.Base, payload []byte) (*Updated, error) {
	const op errors.Op = "bucket.buildUpdated"

	var u Updated

	err := json.Unmarshal(payload, &u)
	if err != nil {
		return nil, errors.New(op, errors.KindUnexpected, "could not unmarshal updated")
	}

	u.Base = eb

	return &u, nil
}

// Use only when decoding serialized events from store
func buildClosed(eb events.Base, payload []byte) (*Closed, error) {
	const op errors.Op = "bucket.buildClosed"

	var c Closed

	err := json.Unmarshal(payload, &c)
	if err != nil {
		return nil, errors.New(op, errors.KindUnexpected, "could not unmarshal closed")
	}

	c.Base = eb

	return &c, nil
}
