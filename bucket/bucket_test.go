package bucket

import (
	"fmt"
	"testing"
	"time"

	"github.com/juelko/bucket/pkg/errors"
	"github.com/juelko/bucket/pkg/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOpenRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		id   ID
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
			args: args{ID("success"), Name("success"), Description("success")},
			want: &OpenRequest{ID("success"), Name("success"), Description("success")},
			err:  nil,
		},
		{
			desc: "invalid name",
			args: args{ID("success"), Name("invalid-name"), Description("success")},
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

			got, err := NewOpenRequest(tC.args.id, tC.args.name, tC.args.desc)

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
				id:   "TestBucket",
				name: "TestName",
				desc: "Test Descritption",
			},
			want: &Opened{
				baseEvent: baseEvent{
					id: "TestBucket",
					v:  1,
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
			assert.True(t, got.OccuredAt().After(begin))

			switch want := tC.want.(type) {
			case *Opened:
				e := got.(*Opened)
				assert.Equal(t, want.StreamID(), e.StreamID())
				assert.Equal(t, want.Version(), e.Version())
				assert.Equal(t, want.Name(), e.Name())
				assert.Equal(t, want.Description(), e.Description())
			case *ErrorEvent:
				e := got.(*ErrorEvent)
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
			args: args{ID("success"), Name("success"), Description("success")},
			want: &UpdateRequest{ID("success"), Name("success"), Description("success")},
			err:  nil,
		},
		{
			desc: "invalid name",
			args: args{ID("success"), Name("invalid-name"), Description("success")},
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

			got, err := NewUpdateRequest(tC.args.id, tC.args.name, tC.args.desc)

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
				name: "NewName",
				desc: "New Descritption",
			},
			stream: openTestStream("TestBucket"),
			want: &Updated{
				baseEvent: baseEvent{
					id: "TestBucket",
					v:  2,
				},
				name: "NewName",
				desc: "New Descritption",
			},
		},
		{
			desc: "closed stream",
			req: &UpdateRequest{
				id:   "TestBucket",
				name: "NewName",
				desc: "New Descritption",
			},
			stream: closedTestStream("TestBucket"),
			want: &ErrorEvent{
				baseEvent: baseEvent{},
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
				name: "NewName",
				desc: "New Descritption",
			},
			stream: closedTestStream("AnotherBucket"),
			want: &ErrorEvent{
				baseEvent: baseEvent{},
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
				name: "NewName",
				desc: "New Descritption",
			},
			stream: []Event{},
			want: &ErrorEvent{
				baseEvent: baseEvent{},
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
			assert.True(t, got.OccuredAt().After(begin))

			switch want := tC.want.(type) {
			case *Updated:
				e := got.(*Updated)
				assert.Equal(t, want.StreamID(), e.StreamID())
				assert.Equal(t, want.Version(), e.Version())
				assert.Equal(t, want.Name(), e.Name())
				assert.Equal(t, want.Description(), e.Description())
			case *ErrorEvent:
				e := got.(*ErrorEvent)
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
		rid request.ID
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
			want: &CloseRequest{ID("success")},
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

			got, err := NewCloseRequest(tC.args.id)

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
				id: "TestBucket",
			},
			stream: openTestStream(ID("TestBucket")),
			want: &Closed{
				baseEvent: baseEvent{
					id: "TestBucket",
					v:  2,
				},
			},
		},
		{
			desc: "closed stream",
			req: &CloseRequest{
				id: "TestBucket",
			},
			stream: closedTestStream(ID("TestBucket")),
			want: &ErrorEvent{
				baseEvent: baseEvent{},
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
				id: "TestBucket",
			},
			stream: closedTestStream(ID("AnotherBucket")),
			want: &ErrorEvent{
				baseEvent: baseEvent{},
				err: &errors.Error{
					Op:    "bucket.stateForClosing",
					Kind:  errors.KindUnexpected,
					Msg:   "Error when building state for closing",
					Wraps: fmt.Errorf("ID Mismatch"),
				},
			},
		},
		{
			desc:   "empty stream",
			req:    &CloseRequest{ID("TestBucket")},
			stream: []Event{},
			want: &ErrorEvent{
				baseEvent: baseEvent{},
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
			assert.True(t, got.OccuredAt().After(begin))

			switch want := tC.want.(type) {
			case *Closed:
				e := got.(*Closed)
				assert.Equal(t, want.StreamID(), e.StreamID())
				assert.Equal(t, want.Version(), e.Version())
			case *ErrorEvent:
				e := got.(*ErrorEvent)
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
			stream: []Event{NewErrorEvent(errors.New("for error event"))},
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

			got := tC.id.Validate()

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

			got := tC.id.Validate()

			assert.Equal(t, tC.want, got)

		})
	}
}

// helper funcs for testing
func openTestStream(id ID) []Event {

	return []Event{
		&Opened{
			baseEvent: baseEvent{
				id: id,
				v:  1,
				at: atOpened,
			},
			name: "TestName",
			desc: "Test Description",
		},
	}
}

func updatedTestStream(id ID) []Event {
	return []Event{
		&Opened{
			baseEvent: baseEvent{
				id: id,
				v:  1,
				at: atOpened,
			},
			name: "TestName",
			desc: "Test Description",
		},
		&Updated{
			baseEvent: baseEvent{
				id: id,
				v:  2,
				at: atUpdated,
			},
			name: "UpdatedName",
			desc: "Updated Description",
		},
	}
}

func closedTestStream(id ID) []Event {
	return []Event{
		&Opened{
			baseEvent: baseEvent{
				id: id,
				v:  1,
				at: atOpened,
			},
			name: "TestName",
			desc: "Test Description",
		},
		&Updated{
			baseEvent: baseEvent{
				id: id,
				v:  2,
				at: atUpdated,
			},
			name: "ClosedName",
			desc: "Closed Description",
		},
		&Closed{
			baseEvent: baseEvent{
				id: id,
				v:  3,
				at: atClosed,
			},
		},
	}
}

func testReqID() request.ID {
	return request.ID("10c0d59e-ca70-46d8-87fb-738be0c9b035")
}

// helper variables for testing
var (
	atOpened  time.Time = time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC)
	atUpdated time.Time = time.Date(2020, time.January, 1, 13, 0, 0, 0, time.UTC)
	atClosed  time.Time = time.Date(2020, time.January, 1, 14, 0, 0, 0, time.UTC)
)
