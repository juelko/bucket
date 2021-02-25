package errors

import (
	"fmt"
	"strings"
)

type Kind int

var kindStrings = []string{"Unkwown", "Unexpected", "Expected", "Validation", "Not Found"}

func (k Kind) String() string {
	if k < 1 && k > 3 {
		return kindStrings[0]
	}
	return kindStrings[k]
}

const (
	KindUnexpected Kind = iota + 1
	KindExpected
	KindValidation
	KindNotFound
)

type Op string

type Error struct {
	Op    Op     // operation
	Kind  Kind   // category of error
	Msg   string // Error meesage
	Wraps error  // Wraps error
}

func (e Error) Error() string {
	var wraps string

	if e.Wraps != nil {
		wraps = e.Wraps.Error()
	}

	var b strings.Builder

	fmt.Fprintf(&b, "operation: %s, kind: %s, error: %s", e.Op, e.Kind, e.Msg)

	if len(wraps) != 0 {
		fmt.Fprintf(&b, ", wraps: %s", wraps)
	}

	return b.String()
}

// Equals reports wheter two Errors has the same content.
// If wrapped errors are Errors type it compares those recursively
// other wrapped errors are not compared
func (e *Error) Equal(t *Error) bool {
	ret := false

	ret = e.Op == t.Op
	ret = e.Kind == t.Kind
	ret = e.Msg == t.Msg

	if wrapped, ok := e.Wraps.(*Error); ok {
		if twrapped, ok := t.Wraps.(*Error); ok {
			wrapped.Equal(twrapped)
		}
	}

	return ret

}

func New(args ...interface{}) *Error {
	e := Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case Kind:
			e.Kind = arg
		case string:
			e.Msg = arg
		case error:
			e.Wraps = arg
		default:

		}
	}
	return &e
}

func Ops(e *Error) []Op {
	res := []Op{e.Op}

	subErr, ok := e.Wraps.(*Error)
	if !ok {
		return res
	}

	res = append(res, Ops(subErr)...)

	return res
}