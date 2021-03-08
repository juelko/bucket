package bucket

import (
	"github.com/juelko/bucket/pkg/errors"
	"github.com/juelko/bucket/pkg/events"
	"github.com/juelko/bucket/pkg/validator"
)

// OpenRequest represent arguments for opening new bucket
type OpenRequest struct {
	ID    events.EntityID
	Title Title
	Desc  Description
}

func (req *OpenRequest) Validate() error {
	const op errors.Op = "bucket.OpenRequest.Validate"

	if err := validator.Validate(req.ID, req.Title); err != nil {
		return errors.New(op, errors.KindValidation, "Invalid arguments", err)
	}

	return nil

}

type CloseRequest struct {
	ID events.EntityID
}

func (req *CloseRequest) Validate() error {
	const op errors.Op = "bucket.CloseRequest.Validate"

	if err := validator.Validate(req.ID); err != nil {
		return errors.New(op, errors.KindValidation, "Invalid arguments", err)
	}

	return nil

}

type UpdateRequest struct {
	ID    events.EntityID
	Title Title
	Desc  Description
}

func (req *UpdateRequest) Validate() error {
	const op errors.Op = "bucket.UpdateRequest.Validate"

	if err := validator.Validate(req.ID, req.Title); err != nil {
		return errors.New(op, errors.KindValidation, "Invalid arguments", err)
	}

	return nil

}

type Reponse struct {
	View *View
	Err  string
}
