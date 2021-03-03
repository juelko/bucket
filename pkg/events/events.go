package events

import (
	"regexp"
	"time"

	"github.com/juelko/bucket/pkg/errors"
	"github.com/juelko/bucket/pkg/request"
	"github.com/juelko/bucket/pkg/validator"
)

type Event interface {
	StreamID() StreamID
	Type() string
	Occured() time.Time
	Version() Version
	RequestID() request.ID
}

func NewBase(id StreamID, v Version, rid request.ID) (Base, error) {
	const op errors.Op = "events.NewBase"

	err := validator.Validate(id, v, rid)
	if err != nil {
		return Base{}, errors.New(op, errors.KindValidation, "Invalid argument", err)
	}

	return Base{
		ID:  id,
		At:  time.Now(),
		V:   v,
		RID: rid,
	}, nil
}

// Base has common fields and methods to all domain events.
// When Base is embedded, Domain Event only needs Type() and Payload() methods to satisfy Event interface
type Base struct {
	ID  StreamID
	RID request.ID
	At  time.Time
	V   Version
}

func (be Base) StreamID() StreamID {
	return be.ID
}

func (be Base) Version() Version {
	return be.V
}

func (be Base) Occured() time.Time {
	return be.At
}

func (be Base) RequestID() request.ID {
	return be.RID
}

// StreamID is identifies for stream of domain events. Use domain entity's identifier as StreamID
type StreamID string

var idRegexp = regexp.MustCompile("^[A-Za-z0-9]{3,64}$")

func (id StreamID) Validate() *errors.Error {
	const op errors.Op = "bucket.ID.Validate"

	if !idRegexp.Match([]byte(id)) {
		return errors.New(op, errors.KindValidation, "Invalid value for ID")
	}

	return nil
}

// Version of the domain entity
type Version uint

func (v Version) Validate() *errors.Error {
	const op errors.Op = "bucket.Version.validate"

	if v == 0 {
		return errors.New(op, errors.KindValidation, "Invalid value for version")
	}

	return nil
}
