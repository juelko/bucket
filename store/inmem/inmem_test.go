package inmem

import (
	"context"
	"testing"

	"github.com/juelko/bucket/bucket"
	"github.com/juelko/bucket/pkg/errors"
	"github.com/juelko/bucket/pkg/events"
	"github.com/stretchr/testify/require"
)

func TestGetStream(t *testing.T) {
	t.Parallel()

	store := NewTestBucketStore()

	testCases := []struct {
		desc string
		args events.EntityID
		want []events.Event
		err  error
	}{
		{
			desc: "happy",
			args: "OpenID",
			want: []events.Event{
				&bucket.Opened{
					Base:       events.Base{ID: "OpenID", V: 1},
					BucketData: bucket.BucketData{Title: "OpenTitle", Description: "Open Description"},
				},
			},
			err: nil,
		},
		{
			desc: "not found",
			args: "NotFoundID",
			want: []events.Event{},
			err:  &errors.Error{Op: "inmem.store.findStream", Kind: 4, Msg: "Stream not found"},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got, err := store.GetStream(context.Background(), tC.args)

			if len(tC.want) != 0 {
				require.Nil(t, err, "error should be nil")
				require.Equal(t, tC.want, got, "events should be equal")
			} else {
				require.Empty(t, got, "event should be empty")
				require.Equal(t, tC.err, err, "errors should be equal")
			}
		})
	}

}

func TestOpenStream(t *testing.T) {
	t.Parallel()

	store := NewTestBucketStore()

	testCases := []struct {
		desc string
		args *bucket.Opened
		want error
	}{
		{
			desc: "happy",
			args: &bucket.Opened{
				Base:       events.Base{ID: "NewID", V: 1},
				BucketData: bucket.BucketData{Title: "NewTitle", Description: "New Description"},
			},
			want: nil,
		},
		{
			desc: "all raedy exists",
			args: &bucket.Opened{
				Base:       events.Base{ID: "OpenID", V: 1},
				BucketData: bucket.BucketData{Title: "OpenTitle", Description: "Open Description"},
			},
			want: &errors.Error{Op: "inmem.store.OpenStream", Kind: 5, Msg: "Allready exists"},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got := store.OpenStream(context.Background(), tC.args)

			require.Equal(t, tC.want, got, "got should be equal")
		})
	}

}
