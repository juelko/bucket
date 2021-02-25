package bucket

import (
	"fmt"
	"time"

	"github.com/juelko/bucket/pkg/errors"

	validator "github.com/go-playground/validator/v10"
)

func newState(id string, stream []Event) (state, error) {

	ret := state{}

	if len(stream) == 0 {
		return ret, errors.New("Empty stream")

	}

	for _, e := range stream {

		if id != e.StreamID() {
			return state{}, errors.New("ID Mismatch")
		}

		switch event := e.(type) {
		case Opened:
			ret.id = event.ID
			ret.name = event.Name
			ret.description = event.Description
			ret.v = event.V
			ret.updated = event.Occured
		case Updated:
			ret.name = event.Name
			ret.description = event.Description
			ret.v = event.V
			ret.updated = event.Occured
		case Closed:
			ret.closed = true
			ret.v = event.V
			ret.updated = event.Occured
		default:
			return state{}, errors.New("Stream contains unkown events")
		}
	}

	return ret, nil
}

type state struct {
	id          string
	name        string
	description string
	closed      bool
	updated     time.Time
	v           uint
}

// Single instance of Validate, because it caches struct info
var validate = validator.New()

func validateStruct(i interface{}) error {
	const op errors.Op = "bucket.validateStruct"

	err := validate.Struct(i)

	if err != nil {
		var ret error

		for _, err := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Invalid value: %v for Field: %s", err.Value(), err.Field())
			ret = errors.New(op, errors.KindValidation, msg, ret)
		}

		return ret
	}

	return nil
}
