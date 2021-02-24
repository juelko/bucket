package bucket

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
)

var ErrInvalidID = errors.New("Invalid Bucket ID")

func NewBucketID() BucketID {
	return BucketID(uuid.NewString())
}

type BucketID string

func (bid BucketID) Validate() error {
	if notValidID(string(bid)) {
		return ErrInvalidID
	}

	return nil
}

func (bid BucketID) String() string {
	return string(bid)
}

var regexpUUIDv4RFC4122 = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$")

func notValidID(s string) bool {
	return !regexpUUIDv4RFC4122.MatchString(s)
}
