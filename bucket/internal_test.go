package bucket

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStruct(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc string
		args interface{}
		want error
	}{
		{
			desc: "BaseRequest ok",
			args: baseRequest{id: "TestBucket", rid: testReqID()},
			want: nil,
		},
		{
			desc: "BaseRequest invalid ID",
			args: baseRequest{id: "Invalid-Bucket-ID", rid: testReqID()},
			want: fmt.Errorf("Invalid value: Invalid-Bucket-ID for Field: ID"),
		},
		{
			desc: "BaseRequest invalid RequestID",
			args: baseRequest{id: "TestBucket", rid: "invalid-request-id"},
			want: fmt.Errorf("Invalid value: invalid-request-id for Field: RequestID"),
		},
		{
			desc: "OpenRequest ok",
			args: &OpenRequest{
				baseRequest: baseRequest{
					id:  "TestBucket",
					rid: testReqID(),
				},
				name:  "TestName",
				dessc: "Test Descritption",
			},
			want: nil,
		},
		{
			desc: "OpenRequest Invalid Name",
			args: &OpenRequest{
				baseRequest: baseRequest{
					id:  "TestBucket",
					rid: testReqID(),
				},
				name:  "Invalid-Name",
				dessc: "Test Descritption",
			},
			want: fmt.Errorf("Invalid value: Invalid-Name for Field: Name"),
		},
		{
			desc: "OpenRequest ok",
			args: &UpdateRequest{
				baseRequest: baseRequest{
					id:  "TestBucket",
					rid: testReqID(),
				},
				name: "TestName",
				desc: "Test Descritption",
			},
			want: nil,
		},
		{
			desc: "UpdateRequest Invalid Name",
			args: &UpdateRequest{
				baseRequest: baseRequest{
					id:  "TestBucket",
					rid: testReqID(),
				},
				name: "Invalid-Name",
				desc: "Test Descritption",
			},
			want: fmt.Errorf("Invalid value: Invalid-Name for Field: Name"),
		},
		{
			desc: "CloseRequest ok",
			args: &CloseRequest{baseRequest{id: "TestBucket", rid: testReqID()}},
			want: nil,
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got := validateStruct(tC.args)

			assert.Equal(t, tC.want, got)

		})
	}
}
