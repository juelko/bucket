package request

import (
	"context"
	"regexp"

	"github.com/google/uuid"
	"github.com/juelko/bucket/pkg/errors"
)

// ID is unique identifier for each request. Format is rfc4122 UUID
type ID string

var ridRegexp = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

func (rid ID) Validate() *errors.Error {
	const op errors.Op = "request.ID.Validate"

	if !ridRegexp.Match([]byte(rid)) {
		return errors.New(op, errors.KindValidation, "Invalid value for request.ID")
	}

	return nil
}

func New() ID {
	return ID(uuid.New().String())
}

func NewContext(ctx context.Context, id ID) context.Context {
	return context.WithValue(ctx, idKey, id)
}

func FromContext(ctx context.Context) (ID, bool) {
	id, ok := ctx.Value(idKey).(ID)
	return id, ok
}

type key int

var idKey key
