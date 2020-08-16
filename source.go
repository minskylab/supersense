package supersense

import "context"

// Source is a new Event emmiter
type Source interface {
	Run(ctx context.Context) error
	Dispose(ctx context.Context)
	Events(ctx context.Context) *chan Event
}
