package bucket

import (
	"context"

	"github.com/juelko/bucket/pkg/events"
)

type Service interface {
	Open(ctx context.Context, req *OpenRequest) (events.Event, error)
	Update(ctx context.Context, req *UpdateRequest) (events.Event, error)
	Close(ctx context.Context, req *CloseRequest) (events.Event, error)
	Get(ctx context.Context, id events.EntityID) (*View, error)
}

type Store interface {
	OpenStream(ctx context.Context, o *Opened) error
	InsertEvent(ctx context.Context, e events.Event) error
	GetStream(ctx context.Context, id events.EntityID) ([]events.Event, error)
}
