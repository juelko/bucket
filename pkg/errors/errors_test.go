package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	t.Parallel()

	second := &Error{Op: "errors.TestError", Kind: KindUnexpected, Msg: "Second Error", Wraps: nil}

	testCases := []struct {
		desc string
		args error
		want string
	}{
		{
			desc: "simple",
			args: &Error{Op: "errors.TestError", Kind: KindUnexpected, Msg: "Test Error", Wraps: nil},
			want: "operation: errors.TestError, kind: Unexpected, error: Test Error",
		},
		{
			desc: "wrapped",
			args: &Error{Op: "errors.TestError", Kind: KindUnexpected, Msg: "Test Error", Wraps: second},
			want: "operation: errors.TestError, kind: Unexpected, error: Test Error",
		},
		{
			desc: "wrapped other",
			args: &Error{Op: "errors.TestError", Kind: KindUnexpected, Msg: "Test Error", Wraps: fmt.Errorf("some other error")},
			want: "operation: errors.TestError, kind: Unexpected, error: Test Error",
		},
	}

	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tC.want, tC.args.Error())
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

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
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			got := New(tC.args.op, tC.args.kind, tC.args.msg, tC.args.wraps)

			assert.Equal(t, tC.want, got)
		})
	}
}

func TestOps(t *testing.T) {
	t.Parallel()

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
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got := Ops(tC.args)

			assert.Equal(t, tC.want, got)
		})
	}
}

func TestEqual(t *testing.T) {

	t.Parallel()

	simple := &Error{Op: "errors.TestEqual", Kind: KindExpected, Msg: "simple", Wraps: nil}

	type args struct {
		in *Error
		to *Error
	}

	testCases := []struct {
		desc string
		args args
		want bool
	}{
		{
			desc: "simple",
			args: args{
				in: simple,
				to: simple,
			},
			want: true,
		},
		{
			desc: "wrapped *Error",
			args: args{
				in: &Error{Op: "errors.TestEqual", Kind: KindExpected, Msg: "wrapped *Error", Wraps: simple},
				to: &Error{Op: "errors.TestEqual", Kind: KindExpected, Msg: "wrapped *Error", Wraps: simple},
			},
			want: true,
		},
		{
			desc: "wrapped same standard error",
			args: args{
				in: &Error{Op: "errors.TestEqual", Kind: KindExpected, Msg: "wrapped same error", Wraps: fmt.Errorf("wrapped")},
				to: &Error{Op: "errors.TestEqual", Kind: KindExpected, Msg: "wrapped same error", Wraps: fmt.Errorf("wrapped")},
			},
			want: true,
		},
		{
			desc: "wrapped different standard error",
			args: args{
				in: &Error{Op: "errors.TestEqual", Kind: KindExpected, Msg: "wrapped same error", Wraps: fmt.Errorf("wrapped")},
				to: &Error{Op: "errors.TestEqual", Kind: KindExpected, Msg: "wrapped same error", Wraps: fmt.Errorf("another wrapped")},
			},
			want: true,
		},
		{
			desc: "different ops",
			args: args{
				in: &Error{Op: "errors.TestEqual", Kind: KindExpected, Msg: "simple", Wraps: nil},
				to: &Error{Op: "errors.TestNotEqual", Kind: KindExpected, Msg: "simple", Wraps: nil},
			},
			want: false,
		},
		{
			desc: "different kind",
			args: args{
				in: &Error{Op: "errors.TestEqual", Kind: KindExpected, Msg: "simple", Wraps: nil},
				to: &Error{Op: "errors.TestEqual", Kind: KindUnexpected, Msg: "simple", Wraps: nil},
			},
			want: false,
		},
		{
			desc: "different msg",
			args: args{
				in: &Error{Op: "errors.TestEqual", Kind: KindExpected, Msg: "simple", Wraps: nil},
				to: &Error{Op: "errors.TestEqual", Kind: KindExpected, Msg: "another", Wraps: nil},
			},
			want: false,
		},
		{
			desc: "different wrapped *Error",
			args: args{
				in: &Error{Op: "errors.TestEqual", Kind: KindExpected, Msg: "simple", Wraps: simple},
				to: &Error{Op: "errors.TestEqual", Kind: KindExpected, Msg: "simple", Wraps: &Error{Op: "errors.TestEqual", Kind: KindExpected, Msg: "another", Wraps: nil}},
			},
			want: false,
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got := tC.args.in.Equal(tC.args.to)

			assert.Equal(t, tC.want, got)

		})
	}
}

func TestKindString(t *testing.T) {

	t.Parallel()

	testCases := []struct {
		desc string
		args Kind
		want string
	}{
		{
			desc: "Unexpected",
			args: KindUnexpected,
			want: "Unexpected",
		},
		{
			desc: "Expected",
			args: KindExpected,
			want: "Expected",
		},
		{
			desc: "Validation",
			args: KindValidation,
			want: "Validation",
		},
		{
			desc: "Not Found",
			args: KindNotFound,
			want: "Not Found",
		},
		{
			desc: "Zero value",
			args: 0,
			want: "Unknown",
		},
		{
			desc: "Out of range",
			args: 100,
			want: "Unknown",
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got := tC.args.String()

			assert.Equal(t, tC.want, got)

		})
	}
}
