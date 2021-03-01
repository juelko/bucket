package validator

import (
	"testing"

	"github.com/juelko/bucket/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc string
		args []Validator
		want error
	}{
		{
			desc: "all ok",
			args: []Validator{newOK(), newOK()},
			want: nil,
		},
		{
			desc: "one err",
			args: []Validator{newError("first")},
			want: &errors.Error{Op: "test.errorValidator.Validate", Kind: errors.KindValidation, Msg: "first", Wraps: error(nil)},
		},
		{
			desc: "two err, one ok",
			args: []Validator{newError("first"), newOK(), newError("second")},
			want: &errors.Error{
				Op:   "test.errorValidator.Validate",
				Kind: errors.KindValidation,
				Msg:  "second",
				Wraps: &errors.Error{
					Op:    "test.errorValidator.Validate",
					Kind:  errors.KindValidation,
					Msg:   "first",
					Wraps: error(nil)}},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got := Validate(tC.args...)

			assert.Equal(t, tC.want, got)

		})
	}
}

// Test helpers
func newOK() Validator {
	return okValidator{}
}

type okValidator struct{}

func (ok okValidator) Validate() *errors.Error {
	return nil
}

func newError(msg string) Validator {
	return errorValidator{msg}
}

type errorValidator struct {
	msg string
}

func (er errorValidator) Validate() *errors.Error {
	const op errors.Op = "test.errorValidator.Validate"

	return errors.New(op, errors.KindValidation, er.msg)
}
