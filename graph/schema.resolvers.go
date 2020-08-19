package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/graph/generated"
)

func (r *queryResolver) Event(ctx context.Context, id string) (*supersense.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *subscriptionResolver) Events(ctx context.Context) (<-chan *supersense.Event, error) {
	pipe := make(chan *supersense.Event, 10)

	go func(eventPipe *chan *supersense.Event) {
		for {
			event := <- r.mux.Events()
			r.mu.Lock()
			// logrus.Warn(event.ID)
			pipe <- &event
			r.mu.Unlock()
		}
	}(&pipe)

	return pipe, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
