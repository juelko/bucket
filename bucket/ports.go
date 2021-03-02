package bucket

import (
	"context"

	"github.com/juelko/bucket/pkg/events"
)

type Service interface {
	Open(ctx context.Context, req *OpenRequest) (*View, error)
	Update(ctx context.Context, req *UpdateRequest) (*View, error)
	Close(ctx context.Context, req *CloseRequest) (*View, error)
	Get(ctx context.Context, id events.StreamID) (*View, error)
}

type Store interface {
	InsertEvent(ctx context.Context, e events.Event) error
	GetStream(ctx context.Context, id events.StreamID) ([]events.Event, error)
}
