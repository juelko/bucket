package bucket

import (
	"regexp"

	"github.com/juelko/bucket/pkg/errors"
)

// ID is the bucket's identifier
type ID string

var idRegexp = regexp.MustCompile("^[A-Za-z0-9]{3,64}$")

func (id ID) Validate() *errors.Error {
	const op errors.Op = "bucket.ID.validate"

	if !idRegexp.Match([]byte(id)) {
		return errors.New(op, errors.KindValidation, "Invalid value for ID")
	}

	return nil
}

// Version of the bucket object
type Version uint

// Name for the bucket
type Name string

var nameRegexp = regexp.MustCompile("^[A-Za-z0-9]{3,64}$")

func (n Name) Validate() *errors.Error {
	const op errors.Op = "bucket.Name.validate"

	if !idRegexp.Match([]byte(n)) {
		return errors.New(op, errors.KindValidation, "Invalid value for Name")
	}
	return nil
}

// Desription for the bucket
type Description string
