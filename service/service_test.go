package bucket

import (
	"testing"

	"github.com/juelko/bucket/store/inmem"
)

func TestServiceOpen(t *testing.T) {

	s := inmem.NewTestBucketStore()

	svc := NewService(s)

	t.Log(svc)

	testCases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}
