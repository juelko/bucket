package bucket

import (
	"context"
	"testing"

	"github.com/juelko/bucket/bucket"
	"github.com/juelko/bucket/pkg/errors"
	"github.com/juelko/bucket/pkg/events"
	"github.com/juelko/bucket/store/inmem"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	t.Parallel()

	svc := NewService(inmem.NewTestBucketStore())

	testCases := []struct {
		desc string
		args *bucket.OpenRequest
		want events.Event
		err  error
	}{
		{
			desc: "happy",
			args: &bucket.OpenRequest{ID: "NewID", Title: "NewTitle", Desc: "New Description"},
			want: &bucket.Opened{Base: events.Base{ID: "NewID", V: 1}, BucketData: bucket.BucketData{Title: "NewTitle", Description: "New Description"}},
			err:  nil,
		},
		{
			desc: "allready exist",
			args: &bucket.OpenRequest{ID: "OpenID", Title: "OpenTitle", Desc: "Open Description"},
			want: nil,
			err: &errors.Error{
				Op:    "inmem.store.OpenStream",
				Kind:  5,
				Msg:   "Allready exists",
				Wraps: nil},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got, err := svc.Open(context.Background(), tC.args)

			if tC.want != nil {
				require.Nil(t, err, "error should be nil")
				require.Equal(t, tC.want, got, "events should be equal")
			} else {
				require.Nil(t, got, "event should be nil")
				require.Equal(t, tC.err, err, "errors should be equal")
			}

		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	svc := NewService(inmem.NewTestBucketStore())

	testCases := []struct {
		desc string
		args *bucket.UpdateRequest
		want events.Event
		err  error
	}{
		{
			desc: "happy",
			args: &bucket.UpdateRequest{ID: "UpdatedID", Title: "NewTitle", Desc: "New Description"},
			want: &bucket.Updated{Base: events.Base{ID: "UpdatedID", V: 3}, BucketData: bucket.BucketData{Title: "NewTitle", Description: "New Description"}},
			err:  nil,
		},
		{
			desc: "not found",
			args: &bucket.UpdateRequest{ID: "NotFoundID", Title: "NewTitle", Desc: "New Description"},
			want: nil,
			err: &errors.Error{
				Op:   "bucket.service.Update",
				Kind: 4, Msg: "Entity not found",
				Wraps: &errors.Error{
					Op:    "inmem.store.GetStream",
					Kind:  4,
					Msg:   "Stream not found",
					Wraps: nil}},
		},
		{
			desc: "closed",
			args: &bucket.UpdateRequest{ID: "ClosedID", Title: "NewTitle", Desc: "New Description"},
			want: nil,
			err: &errors.Error{
				Op:   "bucket.service.Update",
				Kind: 2,
				Msg:  "update not allowed",
				Wraps: &errors.Error{
					Op:    "bucket.updating",
					Kind:  2,
					Msg:   "Bucket is closed",
					Wraps: nil}},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got, err := svc.Update(context.Background(), tC.args)

			if tC.want != nil {
				require.Nil(t, err, "error should be nil")
				require.Equal(t, tC.want, got, "events should be equal")
			} else {
				require.Nil(t, got, "event should be nil")
				require.Equal(t, tC.err, err, "errors should be equal")
			}

		})
	}
}

func TestClose(t *testing.T) {
	t.Parallel()

	svc := NewService(inmem.NewTestBucketStore())

	testCases := []struct {
		desc string
		args *bucket.CloseRequest
		want events.Event
		err  error
	}{
		{
			desc: "happy",
			args: &bucket.CloseRequest{ID: "UpdatedID"},
			want: &bucket.Closed{Base: events.Base{ID: "UpdatedID", V: 3}},
			err:  nil,
		},
		{
			desc: "not found",
			args: &bucket.CloseRequest{ID: "NotFoundID"},
			want: nil,
			err: &errors.Error{
				Op:   "bucket.service.Close",
				Kind: 4, Msg: "Entity not found",
				Wraps: &errors.Error{
					Op:    "inmem.store.GetStream",
					Kind:  4,
					Msg:   "Stream not found",
					Wraps: nil}},
		},
		{
			desc: "closed",
			args: &bucket.CloseRequest{ID: "ClosedID"},
			want: nil,
			err: &errors.Error{
				Op:   "bucket.service.Close",
				Kind: 2,
				Msg:  "closing not allowed",
				Wraps: &errors.Error{
					Op:    "bucket.closing",
					Kind:  2,
					Msg:   "Bucket allready closed",
					Wraps: nil}},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got, err := svc.Close(context.Background(), tC.args)

			if tC.want != nil {
				require.Nil(t, err, "error should be nil")
				require.Equal(t, tC.want, got, "events should be equal")
			} else {
				require.Nil(t, got, "event should be nil")
				require.Equal(t, tC.err, err, "errors should be equal")
			}

		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	svc := NewService(inmem.NewTestBucketStore())

	testCases := []struct {
		desc string
		args events.EntityID
		want *bucket.View
		err  error
	}{
		{
			desc: "happy",
			args: "OpenID",
			want: &bucket.View{
				ID:          "OpenID",
				Title:       "OpenTitle",
				Description: "Open Description",
				Version:     1,
				IsClosed:    false,
			},
			err: nil,
		},
		{
			desc: "not found",
			args: "NotFoundID",
			want: nil,
			err: &errors.Error{
				Op:   "bucket.service.Get",
				Kind: 4,
				Msg:  "Entity not found",
				Wraps: &errors.Error{
					Op:    "inmem.store.GetStream",
					Kind:  4,
					Msg:   "Stream not found",
					Wraps: nil}},
		},
		{
			desc: "invalid id",
			args: "Invalid-ID!",
			want: nil,
			err: &errors.Error{
				Op:   "bucket.service.Get",
				Kind: 3,
				Msg:  "invalid request",
				Wraps: &errors.Error{
					Op:    "bucket.ID.Validate",
					Kind:  3,
					Msg:   "Invalid value for ID",
					Wraps: nil,
				},
			},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got, err := svc.Get(context.Background(), tC.args)

			if tC.want != nil {
				require.Nil(t, err, "error should be nil")
				require.Equal(t, tC.want, got, "events should be equal")
			} else {
				require.Nil(t, got, "event should be nil")
				require.Equal(t, tC.err, err, "errors should be equal")
			}

		})
	}
}
