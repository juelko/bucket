package bucket

import (
	"context"

	"github.com/juelko/bucket/bucket"
	"github.com/juelko/bucket/pkg/errors"
	"github.com/juelko/bucket/pkg/events"
)

func NewService(s bucket.Store) bucket.Service {
	return &service{s}
}

type service struct {
	store bucket.Store
}

func (svc *service) Open(ctx context.Context, req *bucket.OpenRequest) (*bucket.View, error) {
	const op errors.Op = "bucket.service.Open"

	if err := req.Validate(); err != nil {
		return nil, errors.New(op, errors.KindValidation, "invalid request", err)
	}

	o := bucket.Open(req)

	err := svc.store.InsertEvent(ctx, o)
	if err != nil {
		return nil, err
	}

	return bucket.NewView(o.EntityID(), o)
}

func (svc *service) Update(ctx context.Context, req *bucket.UpdateRequest) (*bucket.View, error) {
	const op errors.Op = "bucket.service.Update"

	if err := req.Validate(); err != nil {
		return nil, errors.New(op, errors.KindValidation, "invalid request", err)
	}

	stream, err := svc.store.GetStream(ctx, req.ID)
	if err != nil {
		return nil, errors.New(op, errors.KindExpected, "could not get stream", err)
	}

	u, err := bucket.Update(req, stream)
	if err != nil {
		return nil, errors.New(op, errors.KindExpected, "update not allowed", err)
	}

	err = svc.store.InsertEvent(ctx, u)
	if err != nil {
		return nil, errors.New(op, errors.KindUnexpected, "could not insert", err)
	}

	stream = append(stream, u)

	return bucket.NewView(u.EntityID(), stream...)
}

func (svc *service) Close(ctx context.Context, req *bucket.CloseRequest) (*bucket.View, error) {
	const op errors.Op = "bucket.service.Close"

	if err := req.Validate(); err != nil {
		return nil, errors.New(op, errors.KindValidation, "invalid request", err)
	}

	stream, err := svc.store.GetStream(ctx, req.ID)
	if err != nil {
		return nil, errors.New(op, errors.KindExpected, "could not get stream", err)
	}

	c, err := bucket.Close(req, stream)
	if err != nil {
		return nil, errors.New(op, errors.KindExpected, "closing not allowed", err)
	}

	err = svc.store.InsertEvent(ctx, c)
	if err != nil {
		return nil, errors.New(op, errors.KindUnexpected, "could not insert", err)
	}

	stream = append(stream, c)

	return bucket.NewView(c.EntityID(), stream...)
}

func (svc *service) Get(ctx context.Context, id events.EntityID) (*bucket.View, error) {
	const op errors.Op = "bucket.service.Get"

	if err := id.Validate(); err != nil {
		return nil, errors.New(op, errors.KindValidation, "invalid request", err)
	}

	stream, err := svc.store.GetStream(ctx, id)
	if err != nil {
		return nil, errors.New(op, errors.KindExpected, "could not get stream", err)
	}

	return bucket.NewView(id, stream...)
}
