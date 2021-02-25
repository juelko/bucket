package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {

	second := &Error{Op: "errors.TestError", Kind: KindUnexpected, Msg: "Second Error", Wraps: nil}

	testCases := []struct {
		desc string
		in   error
		want string
	}{
		{
			desc: "simple",
			in:   &Error{Op: "errors.TestError", Kind: KindUnexpected, Msg: "Test Error", Wraps: nil},
			want: "operation: errors.TestError, kind: Unexpected, error: Test Error",
		},
		{
			desc: "wrapped",
			in:   &Error{Op: "errors.TestError", Kind: KindUnexpected, Msg: "Test Error", Wraps: second},
			want: "operation: errors.TestError, kind: Unexpected, error: Test Error, wraps: operation: errors.TestError, kind: Unexpected, error: Second Error",
		},
		{
			desc: "wrapped other",
			in:   &Error{Op: "errors.TestError", Kind: KindUnexpected, Msg: "Test Error", Wraps: fmt.Errorf("some other error")},
			want: "operation: errors.TestError, kind: Unexpected, error: Test Error, wraps: some other error",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tC.want, tC.in.Error())
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		op    Op
		kind  Kind
		msg   string
		wraps error
	}
	testCases := []struct {
		desc string
		args args
		want error
	}{
		{
			desc: "simple",
			args: args{"errors.TestNew", KindExpected, "simple", nil},
			want: &Error{Op: "errors.TestNew", Kind: KindExpected, Msg: "simple", Wraps: nil},
		},
		{
			desc: "wraps",
			args: args{"errors.TestNew", KindExpected, "simple", fmt.Errorf("wrapped")},
			want: &Error{Op: "errors.TestNew", Kind: KindExpected, Msg: "simple", Wraps: fmt.Errorf("wrapped")},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := New(tC.args.op, tC.args.kind, tC.args.msg, tC.args.wraps)

			assert.Equal(t, tC.want, got)
		})
	}
}

func TestOps(t *testing.T) {

	second := &Error{Op: "errors.TestWrapped", Kind: KindExpected, Msg: "second", Wraps: nil}

	testCases := []struct {
		desc string
		args *Error
		want []Op
	}{
		{
			desc: "single",
			args: &Error{Op: "errors.TestOps", Kind: KindExpected, Msg: "simple", Wraps: nil},
			want: []Op{"errors.TestOps"},
		},
		{
			desc: "two",
			args: &Error{Op: "errors.TestOps", Kind: KindExpected, Msg: "simple", Wraps: second},
			want: []Op{"errors.TestOps", "errors.TestWrapped"},
		},
		{
			desc: "wraps other error",
			args: &Error{Op: "errors.TestOps", Kind: KindExpected, Msg: "simple", Wraps: fmt.Errorf("wrapped")},
			want: []Op{"errors.TestOps"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			got := Ops(tC.args)

			assert.Equal(t, tC.want, got)
		})
	}
}
