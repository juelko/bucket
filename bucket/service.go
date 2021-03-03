package bucket

import (
	"context"

	"github.com/juelko/bucket/pkg/errors"
	"github.com/juelko/bucket/pkg/events"
	"github.com/juelko/bucket/pkg/request"
	"github.com/juelko/bucket/pkg/validator"
)

func NewService(s Store) Service {
	return &service{s}
}

type service struct {
	store Store
}

func (svc *service) Open(ctx context.Context, req *OpenRequest) (*View, error) {
	const op errors.Op = "bucket.service.Open"

	if err := req.Validate(); err != nil {
		return nil, errors.New(op, errors.KindValidation, "invalid request", err)
	}

	o := open(req)

	err := svc.store.InsertEvent(ctx, o)
	if err != nil {
		return nil, err
	}

	return NewView(o.StreamID(), o)
}

func (svc *service) Update(ctx context.Context, req *UpdateRequest) (*View, error) {
	const op errors.Op = "bucket.service.Update"

	if err := req.Validate(); err != nil {
		return nil, errors.New(op, errors.KindValidation, "invalid request", err)
	}

	stream, err := svc.store.GetStream(ctx, req.ID)
	if err != nil {
		return nil, errors.New(op, errors.KindExpected, "could not get stream", err)
	}

	u, err := update(req, stream)
	if err != nil {
		return nil, errors.New(op, errors.KindExpected, "update not allowed", err)
	}

	err = svc.store.InsertEvent(ctx, u)
	if err != nil {
		return nil, errors.New(op, errors.KindUnexpected, "could not insert", err)
	}

	stream = append(stream, u)

	return NewView(u.StreamID(), stream...)
}

func (svc *service) Close(ctx context.Context, req *CloseRequest) (*View, error) {
	const op errors.Op = "bucket.service.Close"

	if err := req.Validate(); err != nil {
		return nil, errors.New(op, errors.KindValidation, "invalid request", err)
	}

	stream, err := svc.store.GetStream(ctx, req.ID)
	if err != nil {
		return nil, errors.New(op, errors.KindExpected, "could not get stream", err)
	}

	c, err := close(req, stream)
	if err != nil {
		return nil, errors.New(op, errors.KindExpected, "closing not allowed", err)
	}

	err = svc.store.InsertEvent(ctx, c)
	if err != nil {
		return nil, errors.New(op, errors.KindUnexpected, "could not insert", err)
	}

	stream = append(stream, c)

	return NewView(c.StreamID(), stream...)
}

func (svc *service) Get(ctx context.Context, id events.StreamID) (*View, error) {
	const op errors.Op = "bucket.service.Get"

	if err := id.Validate(); err != nil {
		return nil, errors.New(op, errors.KindValidation, "invalid request", err)
	}

	stream, err := svc.store.GetStream(ctx, id)
	if err != nil {
		return nil, errors.New(op, errors.KindExpected, "could not get stream", err)
	}

	return NewView(id, stream...)
}

// OpenRequest represent arguments for opening new bucket
type OpenRequest struct {
	ID    events.StreamID
	RID   request.ID
	Title Title
	Desc  Description
}

func (req *OpenRequest) Validate() error {
	const op errors.Op = "bucket.OpenRequest.Validate"

	if err := validator.Validate(req.ID, req.RID, req.Title); err != nil {
		return errors.New(op, errors.KindValidation, "Invalid arguments", err)
	}

	return nil

}

type CloseRequest struct {
	ID  events.StreamID
	RID request.ID
}

func (req *CloseRequest) Validate() error {
	const op errors.Op = "bucket.CloseRequest.Validate"

	if err := validator.Validate(req.ID, req.RID); err != nil {
		return errors.New(op, errors.KindValidation, "Invalid arguments", err)
	}

	return nil

}

type UpdateRequest struct {
	ID    events.StreamID
	RID   request.ID
	Title Title
	Desc  Description
}

func (req *UpdateRequest) Validate() error {
	const op errors.Op = "bucket.UpdateRequest.Validate"

	if err := validator.Validate(req.ID, req.RID, req.Title); err != nil {
		return errors.New(op, errors.KindValidation, "Invalid arguments", err)
	}

	return nil

}

type Reponse struct {
	View *View
	Err  string
}
