package bucket

import (
	"context"
)

type Service interface {
	Open(ctx context.Context, req *OpenRequest) (*View, error)
	Update(ctx context.Context, req *UpdateRequest) (*View, error)
	Close(ctx context.Context, req *CloseRequest) (*View, error)
	Find(ctx context.Context, id ID) (*View, error)
}

type Store interface {
	NewStream(ctx context.Context, o Opened) error
	GetStream(ctx context.Context, id ID) []Event
	Insert(ctx context.Context, e Event) error
}
