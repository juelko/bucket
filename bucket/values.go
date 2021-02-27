package bucket

import (
	"regexp"

	"github.com/juelko/bucket/pkg/errors"
)

// ID is the bucket's identifier
type ID string

var idRegexp = regexp.MustCompile("^[A-Za-z0-9]{3-64}$")

func (id ID) validate() error {
	const op errors.Op = "bucket.ID.validate"

	if !idRegexp.Match([]byte(id)) {
		return errors.New(op, errors.KindValidation, "Invalid value for ID")
	}

	return nil
}

// RequestID is unique identifier for each request. Format is rfc4122 uuid v4
type RequestID string

var ridRegexp = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$")

func (rid RequestID) validate() error {
	const op errors.Op = "bucket.RequestID.validate"

	if !idRegexp.Match([]byte(rid)) {
		return errors.New(op, errors.KindValidation, "Invalid value for RequestID")
	}

	return nil
}

// Version of the bucket object
type Version uint

// Name for the bucket
type Name string

var nameRegexp = regexp.MustCompile("^[A-Za-z0-9]{3-64}$")

func (n Name) validate() error {
	const op errors.Op = "bucket.Name.validate"

	if !idRegexp.Match([]byte(n)) {
		return errors.New(op, errors.KindValidation, "Invalid value for Name")
	}
	return nil
}

// Desription for the bucket
type Description string

type validator interface {
	validate() error
}

func validateArgs(args ...validator) error {
	var ret error

	for _, arg := range args {
		ret = arg.validate()
	}

	return ret
}
