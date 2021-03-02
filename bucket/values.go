package bucket

import (
	"regexp"

	"github.com/juelko/bucket/pkg/errors"
)

// Title for the bucket
type Title string

var titleRegexp = regexp.MustCompile(`^[\w -]{3,64}$`)

func (t Title) Validate() *errors.Error {
	const op errors.Op = "bucket.Title.Validate"

	if !titleRegexp.Match([]byte(t)) {
		return errors.New(op, errors.KindValidation, "Invalid value for Title")
	}

	return nil
}

// Desription for the bucket
type Description string
