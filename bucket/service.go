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

	return nil, errors.New("Not implemented")
}

func (svc *service) Update(ctx context.Context, req *UpdateRequest) (*View, error) {

	return nil, errors.New("Not implemented")
}

func (svc *service) Close(ctx context.Context, req *CloseRequest) (*View, error) {

	return nil, errors.New("Not implemented")
}

func (svc *service) Get(ctx context.Context, id events.StreamID) (*View, error) {
	return nil, errors.New("Not implemented")
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
