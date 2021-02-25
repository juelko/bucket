package bucket

import (
	"testing"

	"github.com/juelko/bucket/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testReqID() string {
	return "3a1fc79b-dc53-4f84-a007-4df6f3e54b5f"
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
				BaseRequest: BaseRequest{
					ID:        "TestBucket",
					RequestID: testReqID(),
				},
				Name:        "TestName",
				Description: "Test Descritption",
			},
			want: Opened{
				BaseEvent: BaseEvent{
					ID:    "TestBucket",
					ReqID: testReqID(),
					V:     1,
				},
				Name:        "TestName",
				Description: "Test Descritption",
			},
		},
		{
			desc: "Invalid name",
			args: &OpenRequest{
				BaseRequest: BaseRequest{
					ID:        "TestBucket",
					RequestID: testReqID(),
				},
				Name:        "Invalid-Name",
				Description: "Test Descritption",
			},
			want: ErrorEvent{
				BaseEvent: BaseEvent{
					ReqID: testReqID(),
				},
				Err: &errors.Error{
					Op:   "bucket.Open",
					Kind: errors.KindValidation,
					Msg:  "Request validation failed",
					Wraps: &errors.Error{
						Op:    "bucket.validateStruct",
						Kind:  errors.KindValidation,
						Msg:   "Invalid value: Invalid-Name for Field: Name",
						Wraps: nil,
					},
				},
			},
		},
		{
			desc: "Invalid ID",
			args: &OpenRequest{
				BaseRequest: BaseRequest{
					ID:        "Invalid-Bucket-ID",
					RequestID: testReqID(),
				},
				Name:        "TestName",
				Description: "Test Descritption",
			},
			want: ErrorEvent{
				BaseEvent: BaseEvent{
					ReqID: testReqID(),
				},
				Err: &errors.Error{
					Op:   "bucket.Open",
					Kind: errors.KindValidation,
					Msg:  "Request validation failed",
					Wraps: &errors.Error{
						Op:    "bucket.validateStruct",
						Kind:  errors.KindValidation,
						Msg:   "Invalid value: Invalid-Bucket-ID for Field: ID",
						Wraps: nil,
					},
				},
			},
		},
		{
			desc: "Invalid request id",
			args: &OpenRequest{
				BaseRequest: BaseRequest{
					ID:        "TestBucket",
					RequestID: "invalid-request-id",
				},
				Name:        "TestName",
				Description: "Test Descritption",
			},
			want: ErrorEvent{
				BaseEvent: BaseEvent{
					ReqID: "invalid-request-id",
				},
				Err: &errors.Error{
					Op:   "bucket.Open",
					Kind: errors.KindValidation,
					Msg:  "Request validation failed",
					Wraps: &errors.Error{
						Op:    "bucket.validateStruct",
						Kind:  errors.KindValidation,
						Msg:   "Invalid value: invalid-request-id for Field: RequestID",
						Wraps: nil,
					},
				},
			},
		},
	}

	for i := range testCases {
		tC := testCases[i]

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got := Open(tC.args)

			switch want := tC.want.(type) {
			case Opened:
				e, ok := got.(Opened)
				require.True(t, ok, "Clould not casr Event to Opened")
				assert.Equal(t, want.StreamID(), e.StreamID(), "ID should be Equal")
				assert.Equal(t, want.Type(), e.Type(), "Type should be Equal")
				assert.Equal(t, want.Version(), e.Version(), "Version should be Equal")
				assert.Equal(t, want.RequestID(), e.RequestID(), "RequestId should match")
			case ErrorEvent:
				e, ok := got.(ErrorEvent)
				require.True(t, ok, "Clould not casr Event to ErrorEvent")
				assert.Zero(t, e.StreamID(), "ID should be Zero value")
				assert.Equal(t, want.Type(), e.Type(), "Type should be Equal")
				assert.Zero(t, e.Version(), "Version should be Zero value")
				assert.Equal(t, want.RequestID(), e.RequestID(), "RequestId should match")
				assert.Equal(t, want.Err, e.Err, "Errors should be equal")

			}
		})
	}
}

func openTestStream(id string) []Event {
	return []Event{
		Opened{
			BaseEvent: BaseEvent{
				ID:    id,
				ReqID: testReqID(),
				V:     1,
			},
			Name:        "TestName",
			Description: "Test Descritption",
		},
	}
}

func closedTestStream(id string) []Event {
	return []Event{
		Opened{
			BaseEvent: BaseEvent{
				ID:    id,
				ReqID: testReqID(),
				V:     1,
			},
			Name:        "TestName",
			Description: "Test Descritption",
		},
		Closed{
			BaseEvent: BaseEvent{
				ID:    id,
				ReqID: testReqID(),
				V:     2,
			},
		},
	}
}
