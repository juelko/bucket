package validator

import "github.com/juelko/bucket/pkg/errors"

type Validator interface {
	Validate() *errors.Error
}

func Validate(args ...Validator) error {
	var ret error

	for _, arg := range args {
		if err := arg.Validate(); err != nil {
			err.Wraps = ret
			ret = err
		}
	}

	return ret
}
