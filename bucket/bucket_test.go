package bucket

import (
	"fmt"
	"testing"

	"github.com/juelko/bucket/pkg/errors"
	"github.com/juelko/bucket/pkg/events"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc string
		args *OpenRequest
		want events.Event
		err  error
	}{
		{
			desc: "happy",
			args: &OpenRequest{
				ID:    "TestBucket",
				Title: "TestTitle",
				Desc:  "Test Descritption",
			},
			want: &Opened{
				events.Base{ID: "TestBucket", V: 1},
				BucketData{Title: "TestTitle", Description: "Test Descritption"},
			},
			err: nil,
		},
	}

	for i := range testCases {
		tC := testCases[i]

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got := open(tC.args)

			want := tC.want.(*Opened)
			e, ok := got.(*Opened)
			require.True(t, ok, "failed type casting got to Updated")

			assert.Equal(t, want.Type(), got.Type())
			assert.Equal(t, want.EntityID(), e.EntityID())
			assert.Equal(t, want.EntityVersion(), e.EntityVersion())
			assert.Equal(t, want.Title, e.Title)
			assert.Equal(t, want.Description, e.Description)
		})
	}
}
func TestUpdate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc   string
		req    *UpdateRequest
		stream []events.Event
		want   events.Event
		err    error
	}{
		{
			desc: "happy",
			req: &UpdateRequest{
				ID:    "TestBucket",
				Title: "NewTitle",
				Desc:  "New Descritption",
			},
			stream: openTestStream("TestBucket"),
			want: &Updated{
				events.Base{ID: "TestBucket", V: 2},
				BucketData{"NewTitle", "New Descritption"},
			},
			err: nil,
		},
		{
			desc: "closed stream",
			req: &UpdateRequest{
				ID:    "TestBucket",
				Title: "NewTitle",
				Desc:  "New Descritption",
			},
			stream: closedTestStream("TestBucket"),
			want:   nil,
			err:    &errors.Error{Op: "bucket.updating", Kind: 2, Msg: "Bucket is closed", Wraps: error(nil)},
		},
		{
			desc: "wrong stream",
			req: &UpdateRequest{
				ID:    "TestBucket",
				Title: "NewTitle",
				Desc:  "New Descritption",
			},
			stream: closedTestStream("AnotherBucket"),
			want:   nil,
			err:    &errors.Error{Op: "bucket.stateForUpdating", Kind: 1, Msg: "Error when building state for update", Wraps: fmt.Errorf("ID Mismatch")},
		},
		{
			desc: "empty stream",
			req: &UpdateRequest{
				ID:    "TestBucket",
				Title: "NewTitle",
				Desc:  "New Descritption",
			},
			stream: []events.Event{},
			want:   nil,
			err:    &errors.Error{Op: "bucket.stateForUpdating", Kind: 1, Msg: "Error when building state for update", Wraps: fmt.Errorf("Empty stream")},
		},
	}

	for i := range testCases {
		tC := testCases[i]

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got, err := update(tC.req, tC.stream)

			if tC.want != nil {
				require.Nil(t, err, "require error to be nil")
				want := tC.want.(*Updated)
				e, ok := got.(*Updated)
				require.True(t, ok, "failed type casting got to Updated")

				assert.Equal(t, want.Type(), got.Type())
				assert.Equal(t, want.EntityID(), e.EntityID())
				assert.Equal(t, want.EntityVersion(), e.EntityVersion())
				assert.Equal(t, want.Title, e.Title)
				assert.Equal(t, want.Description, e.Description)
			} else {
				require.Nil(t, got, "require got to be nil")
				assert.Equal(t, tC.err, err)
			}
		})
	}
}

func TestClose(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc   string
		req    *CloseRequest
		stream []events.Event
		want   events.Event
		err    error
	}{
		{
			desc: "happy",
			req: &CloseRequest{
				ID: "TestBucket",
			},
			stream: openTestStream("TestBucket"),
			want:   &Closed{events.Base{ID: "TestBucket", V: 2}},
		},
		{
			desc:   "closed stream",
			req:    &CloseRequest{ID: "TestBucket"},
			stream: closedTestStream("TestBucket"),
			want:   nil,
			err:    &errors.Error{Op: "bucket.closing", Kind: 2, Msg: "Bucket allready closed", Wraps: error(nil)},
		},
		{
			desc:   "wrong stream",
			req:    &CloseRequest{ID: "TestBucket"},
			stream: closedTestStream("AnotherBucket"),
			want:   nil,
			err:    &errors.Error{Op: "bucket.stateForClosing", Kind: 1, Msg: "Error when building state for closing", Wraps: fmt.Errorf("ID Mismatch")},
		},
		{
			desc:   "empty stream",
			req:    &CloseRequest{ID: "TestBucket"},
			stream: []events.Event{},
			want:   nil,
			err:    &errors.Error{Op: "bucket.stateForClosing", Kind: 1, Msg: "Error when building state for closing", Wraps: fmt.Errorf("Empty stream")},
		},
	}

	for i := range testCases {
		tC := testCases[i]

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got, err := close(tC.req, tC.stream)

			if tC.want != nil {
				require.Nil(t, err, "require error to be nil")
				want := tC.want.(*Closed)
				e, ok := got.(*Closed)
				require.True(t, ok, "failed type casting got to Closed")

				assert.Equal(t, want.Type(), e.Type())
				assert.Equal(t, want.EntityID(), e.EntityID())
				assert.Equal(t, want.EntityVersion(), e.EntityVersion())
			} else {
				require.Nil(t, got, "require got to be nil")
				assert.Equal(t, tC.err, err)
			}
		})
	}
}

func TestNewView(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc   string
		id     events.EntityID
		stream []events.Event
		want   *View
		err    error
	}{
		{
			desc:   "ok",
			id:     events.EntityID("TestView"),
			stream: closedTestStream("TestView"),
			want:   &View{ID: "TestView", Title: "ClosedTitle", Description: "Closed Description", Version: 3, IsClosed: true},
			err:    nil,
		},
		{
			desc:   "empty stream",
			id:     events.EntityID("TestView"),
			stream: []events.Event{},
			want:   nil,
			err:    &errors.Error{Op: "bucket.NewView", Kind: errors.KindUnexpected, Msg: "Error when build state for view", Wraps: fmt.Errorf("Empty stream")},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got, err := NewView(tC.id, tC.stream...)

			assert.Equal(t, tC.want, got)
			assert.Equal(t, tC.err, err)

		})
	}
}

// Tests for internals
func TestBuildState(t *testing.T) {
	testCases := []struct {
		desc   string
		id     events.EntityID
		stream []events.Event
		want   state
		err    error
	}{
		{
			desc:   "open",
			id:     events.EntityID("TestView"),
			stream: openTestStream("TestView"),
			want: state{
				id:     events.EntityID("TestView"),
				title:  Title("TestTitle"),
				desc:   Description("Test Description"),
				closed: false,
				v:      1,
			},
			err: nil,
		},
		{
			desc:   "updated",
			id:     events.EntityID("TestView"),
			stream: updatedTestStream("TestView"),
			want: state{
				id:     events.EntityID("TestView"),
				title:  Title("UpdatedTitle"),
				desc:   Description("Updated Description"),
				closed: false,
				v:      2,
			},
			err: nil,
		},
		{
			desc:   "closed",
			id:     events.EntityID("TestView"),
			stream: closedTestStream("TestView"),
			want: state{
				id:     events.EntityID("TestView"),
				title:  Title("ClosedTitle"),
				desc:   Description("Closed Description"),
				closed: true,
				v:      3,
			},
			err: nil,
		},
		{
			desc:   "id mismatch",
			id:     events.EntityID("TestView"),
			stream: openTestStream("DifferentStream"),
			want:   state{},
			err:    fmt.Errorf("ID Mismatch"),
		},
		{
			desc:   "empty stream",
			id:     events.EntityID("TestView"),
			stream: []events.Event{},
			want:   state{},
			err:    fmt.Errorf("Empty stream"),
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

func TestTitleValidation(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc  string
		title Title
		want  *errors.Error
	}{
		{
			desc:  "ok",
			title: Title("Test-Bucket Title_1234"),
			want:  nil,
		},
		{
			desc:  "illegal chars",
			title: Title(`<sript>alert("PWND")</script>`),
			want:  &errors.Error{Op: "bucket.Title.Validate", Kind: errors.KindValidation, Msg: "Invalid value for Title", Wraps: nil},
		},
		{
			desc:  "too short",
			title: Title("ID"),
			want:  &errors.Error{Op: "bucket.Title.Validate", Kind: errors.KindValidation, Msg: "Invalid value for Title", Wraps: nil},
		},
		{
			desc:  "too long",
			title: Title("thisBucketNameisMoreThan64CharactersLong012345678901234567890123456789"),
			want:  &errors.Error{Op: "bucket.Title.Validate", Kind: errors.KindValidation, Msg: "Invalid value for Title", Wraps: nil},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got := tC.title.Validate()

			assert.Equal(t, tC.want, got)

		})
	}
}

// helper funcs for testing
func openTestStream(id events.EntityID) []events.Event {

	return []events.Event{
		&Opened{
			events.Base{ID: id, V: 1},
			BucketData{"TestTitle", "Test Description"},
		},
	}
}

func updatedTestStream(id events.EntityID) []events.Event {
	return []events.Event{
		&Opened{
			events.Base{ID: id, V: 1},
			BucketData{"TestTitle", "Test Description"},
		},
		&Updated{
			events.Base{ID: id, V: 2},
			BucketData{"UpdatedTitle", "Updated Description"},
		},
	}
}

func closedTestStream(id events.EntityID) []events.Event {
	return []events.Event{
		&Opened{
			events.Base{ID: id, V: 1},
			BucketData{"TestTitle", "Test Description"},
		},
		&Updated{
			events.Base{ID: id, V: 2},
			BucketData{"ClosedTitle", "Closed Description"},
		},
		&Closed{events.Base{ID: id, V: 3}},
	}
}
