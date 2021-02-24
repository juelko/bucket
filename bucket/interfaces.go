package bucket

import (
	"context"
)

type Service interface {
	Open(ctx context.Context, req OpenRequest) Event
	Update(ctx context.Context, req UpdateRequest) Event
	Close(ctx context.Context, req CloseRequest) Event
	Find(ctx context.Context, id string) View
}

type Store interface {
	GetStream(ctx context.Context, id string) []Event
	Insert(ctx context.Context, e Event) error
}
