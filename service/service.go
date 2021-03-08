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

func (svc *service) Open(ctx context.Context, req *bucket.OpenRequest) (events.Event, error) {

	o, err := bucket.Open(req)
	if err != nil {
		return nil, err
	}

	if err := svc.store.OpenStream(ctx, o); err != nil {
		return nil, err
	}

	return o, nil
}

func (svc *service) Update(ctx context.Context, req *bucket.UpdateRequest) (events.Event, error) {
	const op errors.Op = "bucket.service.Update"

	if err := req.Validate(); err != nil {
		return nil, errors.New(op, errors.KindValidation, "invalid request", err)
	}

	stream, err := svc.store.GetStream(ctx, req.ID)
	if err != nil {
		return nil, errors.New(op, errors.KindNotFound, "Entity not found", err)
	}

	u, err := bucket.Update(req, stream)
	if err != nil {
		return nil, errors.New(op, errors.KindExpected, "update not allowed", err)
	}

	err = svc.store.InsertEvent(ctx, u)
	if err != nil {
		return nil, errors.New(op, errors.KindUnexpected, "could not insert", err)
	}

	return u, nil
}

func (svc *service) Close(ctx context.Context, req *bucket.CloseRequest) (events.Event, error) {
	const op errors.Op = "bucket.service.Close"

	if err := req.Validate(); err != nil {
		return nil, errors.New(op, errors.KindValidation, "invalid request", err)
	}

	stream, err := svc.store.GetStream(ctx, req.ID)
	if err != nil {
		return nil, errors.New(op, errors.KindNotFound, "Entity not found", err)
	}

	c, err := bucket.Close(req, stream)
	if err != nil {
		return nil, errors.New(op, errors.KindExpected, "closing not allowed", err)
	}

	err = svc.store.InsertEvent(ctx, c)
	if err != nil {
		return nil, errors.New(op, errors.KindUnexpected, "could not insert", err)
	}

	return c, nil
}

func (svc *service) Get(ctx context.Context, id events.EntityID) (*bucket.View, error) {
	const op errors.Op = "bucket.service.Get"

	if err := id.Validate(); err != nil {
		return nil, errors.New(op, errors.KindValidation, "invalid request", err)
	}

	stream, err := svc.store.GetStream(ctx, id)
	if err != nil {
		return nil, errors.New(op, errors.KindNotFound, "Entity not found", err)
	}

	return bucket.NewView(id, stream...)
}
