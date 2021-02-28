package bucket

import (
	"fmt"
	"testing"
	"time"

	"github.com/juelko/bucket/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc string
		args *OpenRequest
		want Event
	}{
		{
			desc: "happy",
			args: &OpenRequest{
				id:    "TestBucket",
				rid:   testReqID(),
				name:  "TestName",
				dessc: "Test Descritption",
			},
			want: Opened{
				baseEvent: baseEvent{
					id:  "TestBucket",
					rid: testReqID(),
					v:   1,
				},
				name: "TestName",
				desc: "Test Descritption",
			},
		},
		{
			desc: "Invalid name",
			args: &OpenRequest{
				id:    "TestBucket",
				rid:   testReqID(),
				name:  "Invalid-Name",
				dessc: "Test Descritption",
			},
			want: ErrorEvent{
				baseEvent: baseEvent{
					rid: testReqID(),
				},
				err: &errors.Error{
					Op:    "bucket.validateOpenRequest",
					Kind:  errors.KindValidation,
					Msg:   "Request validation failed",
					Wraps: fmt.Errorf("Invalid value: Invalid-Name for Field: Name"),
				},
			},
		},
	}

	for i := range testCases {
		tC := testCases[i]

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			begin := time.Now()

			got := Open(tC.args)

			require.Equal(t, tC.want.Type(), got.Type())
			assert.Equal(t, tC.want.RequestID(), got.RequestID())
			assert.True(t, got.OccuredAt().After(begin))

			switch want := tC.want.(type) {
			case Opened:
				e := got.(Opened)
				assert.Equal(t, want.StreamID(), e.StreamID())
				assert.Equal(t, want.Version(), e.Version())
				assert.Equal(t, want.name, e.name)
				assert.Equal(t, want.desc, e.desc)
			case ErrorEvent:
				e := got.(ErrorEvent)
				assert.Zero(t, e.StreamID(), "ID should be Zero value")
				assert.Zero(t, e.Version(), "Version should be Zero value")
				assert.Equal(t, want.Error(), e.Error(), "Errors should be equal")
			default:
				assert.Fail(t, "Unexpedted Event type")
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc   string
		req    *UpdateRequest
		stream []Event
		want   Event
	}{
		{
			desc: "happy",
			req: &UpdateRequest{
				id:   "TestBucket",
				rid:  testReqID(),
				name: "NewName",
				desc: "New Descritption",
			},
			stream: openTestStream("TestBucket"),
			want: Updated{
				baseEvent: baseEvent{
					id:  "TestBucket",
					rid: testReqID(),
					v:   2,
				},
				name: "NewName",
				desc: "New Descritption",
			},
		},
		{
			desc: "invalid new name",
			req: &UpdateRequest{
				id:   "TestBucket",
				rid:  testReqID(),
				name: "Invalid-New-Name",
				desc: "New Descritption",
			},
			stream: closedTestStream("TestBucket"),
			want: ErrorEvent{
				baseEvent: baseEvent{
					rid: testReqID(),
				},
				err: &errors.Error{
					Op:    "bucket.validateUpdateRequest",
					Kind:  errors.KindValidation,
					Msg:   "Request validation failed",
					Wraps: fmt.Errorf("Invalid value: Invalid-New-Name for Field: Name"),
				},
			},
		},
		{
			desc: "closed stream",
			req: &UpdateRequest{
				id:   "TestBucket",
				rid:  testReqID(),
				name: "NewName",
				desc: "New Descritption",
			},
			stream: closedTestStream("TestBucket"),
			want: ErrorEvent{
				baseEvent: baseEvent{
					rid: testReqID(),
				},
				err: &errors.Error{
					Op:    "bucket.updating",
					Kind:  errors.KindExpected,
					Msg:   "Bucket is closed",
					Wraps: nil,
				},
			},
		},
		{
			desc: "wrong stream",
			req: &UpdateRequest{
				id:   "TestBucket",
				rid:  testReqID(),
				name: "NewName",
				desc: "New Descritption",
			},
			stream: closedTestStream("AnotherBucket"),
			want: ErrorEvent{
				baseEvent: baseEvent{
					rid: testReqID(),
				},
				err: &errors.Error{
					Op:    "bucket.buildStateForUpdate",
					Kind:  errors.KindUnexpected,
					Msg:   "Error when building state for update",
					Wraps: fmt.Errorf("ID Mismatch"),
				},
			},
		},
		{
			desc: "empty stream",
			req: &UpdateRequest{
				id:   "TestBucket",
				rid:  testReqID(),
				name: "NewName",
				desc: "New Descritption",
			},
			stream: []Event{},
			want: ErrorEvent{
				baseEvent: baseEvent{
					rid: testReqID(),
				},
				err: &errors.Error{
					Op:    "bucket.buildStateForUpdate",
					Kind:  errors.KindUnexpected,
					Msg:   "Error when building state for update",
					Wraps: fmt.Errorf("Empty stream"),
				},
			},
		},
	}

	for i := range testCases {
		tC := testCases[i]

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			begin := time.Now()

			got := Update(tC.req, tC.stream)

			require.Equal(t, tC.want.Type(), got.Type())
			assert.Equal(t, tC.want.RequestID(), got.RequestID())
			assert.True(t, got.OccuredAt().After(begin))

			switch want := tC.want.(type) {
			case Updated:
				e := got.(Updated)
				assert.Equal(t, want.StreamID(), e.StreamID())
				assert.Equal(t, want.Version(), e.Version())
				assert.Equal(t, want.name, e.name)
				assert.Equal(t, want.desc, e.desc)
			case ErrorEvent:
				e := got.(ErrorEvent)
				assert.Zero(t, e.StreamID(), "ID should be Zero value")
				assert.Zero(t, e.Version(), "Version should be Zero value")
				assert.Equal(t, want.Error(), e.Error(), "Errors should be equal")
			default:
				assert.Fail(t, "Unexpedted Event type")
			}
		})
	}
}

func TestClose(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc   string
		req    *CloseRequest
		stream []Event
		want   Event
	}{
		{
			desc: "happy",
			req: &CloseRequest{
				id:  "TestBucket",
				rid: testReqID(),
			},
			stream: openTestStream(ID("TestBucket")),
			want: Closed{
				baseEvent: baseEvent{
					id:  "TestBucket",
					rid: testReqID(),
					v:   2,
				},
			},
		},
		{
			desc: "closed stream",
			req: &CloseRequest{
				id:  "TestBucket",
				rid: testReqID(),
			},
			stream: closedTestStream(ID("TestBucket")),
			want: ErrorEvent{
				baseEvent: baseEvent{
					rid: testReqID(),
				},
				err: &errors.Error{
					Op:    "bucket.closing",
					Kind:  errors.KindExpected,
					Msg:   "Bucket allready closed",
					Wraps: nil,
				},
			},
		},
		{
			desc: "wrong stream",
			req: &CloseRequest{
				id:  "TestBucket",
				rid: testReqID(),
			},
			stream: closedTestStream(ID("AnotherBucket")),
			want: ErrorEvent{
				baseEvent: baseEvent{
					rid: testReqID(),
				},
				err: &errors.Error{
					Op:    "bucket.stateForClosing",
					Kind:  errors.KindUnexpected,
					Msg:   "Error when building state for closing",
					Wraps: fmt.Errorf("ID Mismatch"),
				},
			},
		},
		{
			desc: "empty stream",
			req: &CloseRequest{
				id:  "TestBucket",
				rid: testReqID(),
			},
			stream: []Event{},
			want: ErrorEvent{
				baseEvent: baseEvent{
					rid: testReqID(),
				},
				err: &errors.Error{
					Op:    "bucket.stateForClosing",
					Kind:  errors.KindUnexpected,
					Msg:   "Error when building state for closing",
					Wraps: fmt.Errorf("Empty stream"),
				},
			},
		},
	}

	for i := range testCases {
		tC := testCases[i]

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			begin := time.Now()

			got := Close(tC.req, tC.stream)

			require.Equal(t, tC.want.Type(), got.Type())
			assert.Equal(t, tC.want.RequestID(), got.RequestID())
			assert.True(t, got.OccuredAt().After(begin))

			switch want := tC.want.(type) {
			case Closed:
				e := got.(Closed)
				assert.Equal(t, want.StreamID(), e.StreamID())
				assert.Equal(t, want.Version(), e.Version())
			case ErrorEvent:
				e := got.(ErrorEvent)
				assert.Zero(t, e.StreamID(), "ID should be Zero value")
				assert.Zero(t, e.Version(), "Version should be Zero value")
				assert.Equal(t, want.Error(), e.Error(), "Errors should be equal")
			default:
				assert.Fail(t, "Unexpedted Event type")
			}
		})
	}
}

func openTestStream(id ID) []Event {
	return []Event{
		Opened{
			baseEvent: baseEvent{
				id:  id,
				rid: testReqID(),
				v:   1,
			},
			name: "TestName",
			desc: "Test Descritption",
		},
	}
}

func closedTestStream(id ID) []Event {
	return []Event{
		Opened{
			baseEvent: baseEvent{
				id:  id,
				rid: testReqID(),
				v:   1,
			},
			name: "TestName",
			desc: "Test Descritption",
		},
		Closed{
			baseEvent: baseEvent{
				id:  id,
				rid: testReqID(),
				v:   2,
			},
		},
	}
}

func testReqID() RequestID {
	return RequestID("3a1fc79b-dc53-4f84-a007-4df6f3e54b5f")
}
