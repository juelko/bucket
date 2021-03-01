package request

import (
	"context"
	"testing"

	"github.com/juelko/bucket/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDValidation(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc string
		id   ID
		want *errors.Error
	}{
		{
			desc: "ok",
			id:   testReqID(),
			want: nil,
		},
		{
			desc: "invalid RequestID",
			id:   ID("this-in-invalid-request-id"),
			want: &errors.Error{Op: "request.ID.Validate", Kind: errors.KindValidation, Msg: "Invalid value for request.ID", Wraps: error(nil)},
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got := tC.id.Validate()

			assert.Equal(t, tC.want, got)

			t.Log(tC.id)

		})
	}
}

func TestContext(t *testing.T) {
	t.Parallel()

	ctx := NewContext(context.Background(), testReqID())

	id, ok := FromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, testReqID(), id)

}

func TestNew(t *testing.T) {
	t.Parallel()

	assert.Nil(t, New().Validate())
}

func testReqID() ID {
	return ID("10c0d59e-ca70-46d8-87fb-738be0c9b035")
}
