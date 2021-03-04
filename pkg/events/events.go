package events

import (
	"regexp"

	"github.com/juelko/bucket/pkg/errors"
	"github.com/juelko/bucket/pkg/validator"
)

type Event interface {
	EntityID() EntityID
	EntityVersion() EntityVersion
	Type() string
	Data() interface{}
}

func NewBase(id EntityID, v EntityVersion) (Base, error) {
	const op errors.Op = "events.NewBase"

	err := validator.Validate(id, v)
	if err != nil {
		return Base{}, errors.New(op, errors.KindValidation, "Invalid argument", err)
	}

	return Base{
		ID: id,
		V:  v,
	}, nil
}

// Base has common fields and methods to all domain events.
// When Base is embedded, Domain Event needs Type() and Data() methods to satisfy Event interface
type Base struct {
	ID EntityID
	V  EntityVersion
}

func (be Base) EntityID() EntityID {
	return be.ID
}

func (be Base) EntityVersion() EntityVersion {
	return be.V
}

// EntityID is identifies for stream of domain events. Use domain entity's identifier as EntityID
type EntityID string

var idRegexp = regexp.MustCompile("^[A-Za-z0-9]{3,64}$")

func (id EntityID) Validate() *errors.Error {
	const op errors.Op = "bucket.ID.Validate"

	if !idRegexp.Match([]byte(id)) {
		return errors.New(op, errors.KindValidation, "Invalid value for ID")
	}

	return nil
}

// EntityVersion of the domain entity
type EntityVersion uint

func (v EntityVersion) Validate() *errors.Error {
	const op errors.Op = "bucket.Version.validate"

	if v == 0 {
		return errors.New(op, errors.KindValidation, "Invalid value for version")
	}

	return nil
}
