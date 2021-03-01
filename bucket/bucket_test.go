package bucket

import (
	"fmt"
	"testing"
	"time"

	"github.com/juelko/bucket/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOpenRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		id   ID
		rid  RequestID
		name Name
		desc Description
	}

	testCases := []struct {
		desc string
		args args
		want *OpenRequest
		err  error
	}{
		{
			desc: "success",
			args: args{ID("success"), testReqID(), Name("success"), Description("success")},
			want: &OpenRequest{ID("success"), testReqID(), Name("success"), Description("success")},
			err:  nil,
		},
		{
			desc: "invalid name",
			args: args{ID("success"), testReqID(), Name("invalid-name"), Description("success")},
			want: nil,
			err: &errors.Error{
				Op:    "bucket.NewOpenRequest",
				Kind:  errors.KindValidation,
				Msg:   "Invalid arguments",
				Wraps: &errors.Error{Op: "bucket.Name.validate", Kind: errors.KindValidation, Msg: "Invalid value for Name", Wraps: error(nil)}},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got, err := NewOpenRequest(tC.args.id, tC.args.rid, tC.args.name, tC.args.desc)

			assert.Equal(t, tC.want, got, "Request should be equal")
			assert.Equal(t, tC.err, err, "error should be equal")

		})
	}
}

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

func TestNewUpdateRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		id   ID
		rid  RequestID
		name Name
		desc Description
	}

	testCases := []struct {
		desc string
		args args
		want *UpdateRequest
		err  error
	}{
		{
			desc: "success",
			args: args{ID("success"), testReqID(), Name("success"), Description("success")},
			want: &UpdateRequest{ID("success"), testReqID(), Name("success"), Description("success")},
			err:  nil,
		},
		{
			desc: "invalid name",
			args: args{ID("success"), testReqID(), Name("invalid-name"), Description("success")},
			want: nil,
			err: &errors.Error{
				Op:    "bucket.NewUpdateRequest",
				Kind:  errors.KindValidation,
				Msg:   "Invalid arguments",
				Wraps: &errors.Error{Op: "bucket.Name.validate", Kind: errors.KindValidation, Msg: "Invalid value for Name", Wraps: error(nil)}},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got, err := NewUpdateRequest(tC.args.id, tC.args.rid, tC.args.name, tC.args.desc)

			assert.Equal(t, tC.want, got, "Request should be equal")
			assert.Equal(t, tC.err, err, "error should be equal")

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
					Op:    "bucket.stateForUpdating",
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
					Op:    "bucket.stateForUpdating",
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

func TestNewCloseRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		id  ID
		rid RequestID
	}

	testCases := []struct {
		desc string
		args args
		want *CloseRequest
		err  error
	}{
		{
			desc: "success",
			args: args{ID("success"), testReqID()},
			want: &CloseRequest{ID("success"), testReqID()},
			err:  nil,
		},
		{
			desc: "invalid ID",
			args: args{ID("InvalidID!"), testReqID()},
			want: nil,
			err: &errors.Error{
				Op:    "bucket.NewCloseRequest",
				Kind:  errors.KindValidation,
				Msg:   "Invalid arguments",
				Wraps: &errors.Error{Op: "bucket.ID.validate", Kind: errors.KindValidation, Msg: "Invalid value for ID", Wraps: error(nil)}},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got, err := NewCloseRequest(tC.args.id, tC.args.rid)

			assert.Equal(t, tC.want, got, "Request should be equal")
			assert.Equal(t, tC.err, err, "error should be equal")

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

func TestNewView(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc   string
		id     ID
		stream []Event
		want   *View
		err    error
	}{
		{
			desc:   "ok",
			id:     ID("TestView"),
			stream: closedTestStream(ID("TestView")),
			want:   &View{ID: "TestView", Name: "ClosedName", Description: "Closed Description", Version: 3, IsClosed: true, LastUpdate: atClosed},
			err:    nil,
		},
		{
			desc:   "empty stream",
			id:     ID("TestView"),
			stream: []Event{},
			want:   nil,
			err:    &errors.Error{Op: "bucket.NewView", Kind: errors.KindUnexpected, Msg: "Error when build state for view", Wraps: fmt.Errorf("Empty stream")},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got, err := NewView(tC.id, tC.stream)

			assert.Equal(t, tC.want, got)
			assert.Equal(t, tC.err, err)

		})
	}
}

// Tests for internals
func TestBuildState(t *testing.T) {
	testCases := []struct {
		desc   string
		id     ID
		stream []Event
		want   state
		err    error
	}{
		{
			desc:   "open",
			id:     ID("OpenStream"),
			stream: openTestStream("OpenStream"),
			want: state{
				id:     ID("OpenStream"),
				name:   Name("TestName"),
				desc:   Description("Test Description"),
				closed: false,
				v:      1,
				last:   atOpened,
			},
			err: nil,
		},
		{
			desc:   "updated",
			id:     ID("UpdatedStream"),
			stream: updatedTestStream("UpdatedStream"),
			want: state{
				id:     ID("UpdatedStream"),
				name:   Name("UpdatedName"),
				desc:   Description("Updated Description"),
				closed: false,
				v:      2,
				last:   atUpdated,
			},
			err: nil,
		},
		{
			desc:   "closed",
			id:     ID("ClosedStream"),
			stream: closedTestStream("ClosedStream"),
			want: state{
				id:     ID("ClosedStream"),
				name:   Name("ClosedName"),
				desc:   Description("Closed Description"),
				closed: true,
				v:      3,
				last:   atClosed,
			},
			err: nil,
		},
		{
			desc:   "id mismatch",
			id:     ID("OpenStream"),
			stream: openTestStream("DifferentStream"),
			want:   state{},
			err:    fmt.Errorf("ID Mismatch"),
		},
		{
			desc:   "empty stream",
			id:     ID("EmptyStream"),
			stream: []Event{},
			want:   state{},
			err:    fmt.Errorf("Empty stream"),
		},
		{
			desc:   "unkonwn event",
			id:     ID(""),
			stream: []Event{newErrorEvent(errors.New("for error event"), testReqID())},
			want:   state{},
			err:    fmt.Errorf("Stream contains unkown events"),
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got, err := buildState(tC.id, tC.stream)

			assert.Equal(t, tC.want, got, "state should be equal")
			assert.Equal(t, tC.err, err, "state should be equal")

		})
	}
}

func TestIDValidation(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc string
		id   ID
		want *errors.Error
	}{
		{
			desc: "ok",
			id:   ID("TestBucketID"),
			want: nil,
		},
		{
			desc: "illegal chars",
			id:   ID("TestBucketID!"),
			want: &errors.Error{Op: "bucket.ID.validate", Kind: errors.KindValidation, Msg: "Invalid value for ID", Wraps: error(nil)},
		},
		{
			desc: "too short",
			id:   ID("ID"),
			want: &errors.Error{Op: "bucket.ID.validate", Kind: errors.KindValidation, Msg: "Invalid value for ID", Wraps: error(nil)},
		},
		{
			desc: "too long",
			id:   ID("thisBucketIDisMoreThan64CharactersLong012345678901234567890123456789"),
			want: &errors.Error{Op: "bucket.ID.validate", Kind: errors.KindValidation, Msg: "Invalid value for ID", Wraps: error(nil)},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got := tC.id.validate()

			assert.Equal(t, tC.want, got)

		})
	}
}

func TestNameValidation(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc string
		id   Name
		want *errors.Error
	}{
		{
			desc: "ok",
			id:   Name("TestBucketName"),
			want: nil,
		},
		{
			desc: "illegal chars",
			id:   Name("TestBucketName!"),
			want: &errors.Error{Op: "bucket.Name.validate", Kind: errors.KindValidation, Msg: "Invalid value for Name", Wraps: error(nil)},
		},
		{
			desc: "too short",
			id:   Name("ID"),
			want: &errors.Error{Op: "bucket.Name.validate", Kind: errors.KindValidation, Msg: "Invalid value for Name", Wraps: error(nil)},
		},
		{
			desc: "too long",
			id:   Name("thisBucketNameisMoreThan64CharactersLong012345678901234567890123456789"),
			want: &errors.Error{Op: "bucket.Name.validate", Kind: errors.KindValidation, Msg: "Invalid value for Name", Wraps: error(nil)},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got := tC.id.validate()

			assert.Equal(t, tC.want, got)

		})
	}
}

func TestRequestIDValidation(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc string
		id   RequestID
		want *errors.Error
	}{
		{
			desc: "ok",
			id:   testReqID(),
			want: nil,
		},
		{
			desc: "invalid RequestID",
			id:   RequestID("this-in-invalid-request-id"),
			want: &errors.Error{Op: "bucket.RequestID.validate", Kind: errors.KindValidation, Msg: "Invalid value for RequestID", Wraps: error(nil)},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got := tC.id.validate()

			assert.Equal(t, tC.want, got)

			t.Log(tC.id)

		})
	}
}

func TestValidateArgs(t *testing.T) {
	t.Parallel()

	type args struct {
		id   ID
		rid  RequestID
		name Name
	}

	testCases := []struct {
		desc string
		args args
		want error
	}{
		{
			desc: "all ok",
			args: args{ID("TestID"), testReqID(), Name("TestName")},
			want: nil,
		},
		{
			desc: "invalid id",
			args: args{ID("InvalidID!"), testReqID(), Name("TestName")},
			want: &errors.Error{Op: "bucket.ID.validate", Kind: errors.KindValidation, Msg: "Invalid value for ID", Wraps: nil},
		},
		{
			desc: "invalid RequestID",
			args: args{ID("TestID"), RequestID("invalid-request-id"), Name("TestName")},
			want: &errors.Error{Op: "bucket.RequestID.validate", Kind: errors.KindValidation, Msg: "Invalid value for RequestID", Wraps: nil},
		},
		{
			desc: "invalid Name",
			args: args{ID("TestID"), testReqID(), Name("InvalidName!")},
			want: &errors.Error{Op: "bucket.Name.validate", Kind: errors.KindValidation, Msg: "Invalid value for Name", Wraps: nil},
		},
		{
			desc: "all invalid",
			args: args{ID("InvalidID!"), RequestID("invalid-request-id"), Name("InvalidName!")},
			want: &errors.Error{
				Op:   "bucket.ID.validate",
				Kind: errors.KindValidation,
				Msg:  "Invalid value for ID",
				Wraps: &errors.Error{
					Op:   "bucket.RequestID.validate",
					Kind: errors.KindValidation,
					Msg:  "Invalid value for RequestID",
					Wraps: &errors.Error{
						Op:    "bucket.Name.validate",
						Kind:  errors.KindValidation,
						Msg:   "Invalid value for Name",
						Wraps: nil}}},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got := validateArgs(tC.args.name, tC.args.rid, tC.args.id)

			assert.Equal(t, tC.want, got)

		})
	}
}

// helper funcs for testing
func openTestStream(id ID) []Event {

	return []Event{
		Opened{
			baseEvent: baseEvent{
				id:  id,
				rid: testReqID(),
				v:   1,
				at:  atOpened,
			},
			name: "TestName",
			desc: "Test Description",
		},
	}
}

func updatedTestStream(id ID) []Event {
	return []Event{
		Opened{
			baseEvent: baseEvent{
				id:  id,
				rid: testReqID(),
				v:   1,
				at:  atOpened,
			},
			name: "TestName",
			desc: "Test Description",
		},
		Updated{
			baseEvent: baseEvent{
				id:  id,
				rid: testReqID(),
				v:   2,
				at:  atUpdated,
			},
			name: "UpdatedName",
			desc: "Updated Description",
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
				at:  atOpened,
			},
			name: "TestName",
			desc: "Test Description",
		},
		Updated{
			baseEvent: baseEvent{
				id:  id,
				rid: testReqID(),
				v:   2,
				at:  atUpdated,
			},
			name: "ClosedName",
			desc: "Closed Description",
		},
		Closed{
			baseEvent: baseEvent{
				id:  id,
				rid: testReqID(),
				v:   3,
				at:  atClosed,
			},
		},
	}
}

func testReqID() RequestID {
	return RequestID("10c0d59e-ca70-46d8-87fb-738be0c9b035")
}

// helper variables for testing
var (
	atOpened  time.Time = time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC)
	atUpdated time.Time = time.Date(2020, time.January, 1, 13, 0, 0, 0, time.UTC)
	atClosed  time.Time = time.Date(2020, time.January, 1, 14, 0, 0, 0, time.UTC)
)
