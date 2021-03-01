package service

import (
	"context"
	"errors"

	"github.com/juelko/bucket/bucket"
)

func NewService(s bucket.Store) bucket.Service {
	return &service{s}
}

type service struct {
	store bucket.Store
}

func (svc *service) Open(ctx context.Context, req *bucket.OpenRequest) (*bucket.View, error) {

	return nil, errors.New("Not implemented")
}

func (svc *service) Update(ctx context.Context, req *bucket.UpdateRequest) (*bucket.View, error) {

	return nil, errors.New("Not implemented")
}

func (svc *service) Close(ctx context.Context, req *bucket.CloseRequest) (*bucket.View, error) {

	return nil, errors.New("Not implemented")
}

func (svc *service) Find(ctx context.Context, id bucket.ID) (*bucket.View, error) {
	return nil, errors.New("Not implemented")
}
