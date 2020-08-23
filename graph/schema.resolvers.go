package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/graph/generated"
	"github.com/minskylab/supersense/graph/model"
	"github.com/pkg/errors"
)

func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*model.AuthResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Broadcast(ctx context.Context, draft model.EventDraft) (string, error) {
	entities := supersense.Entities{
		Tags:  []string{},
		Media: []supersense.MediaEntity{},
		Urls:  []supersense.URLEntity{},
	}

	url := ""
	if draft.ShareURL != nil {
		url = *draft.ShareURL
	}


	if draft.Entities != nil {
		for _, tag := range draft.Entities.Tags {
			entities.Tags = append(entities.Tags, tag)
		}
	}

	r.spokesman.Broadcast(draft.Title, draft.Message, entities, url, nil)

	return draft.Message, nil
}

func (r *queryResolver) SharedBoard(ctx context.Context, buffer int) ([]*supersense.Event, error) {
	if buffer < 1 || buffer > 100 { // hard coded limit
		buffer = 100
	}

	currentState, err := r.store.CurrentSharedState(int64(buffer))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return currentState.Board, nil
}

func (r *subscriptionResolver) EventStream(ctx context.Context, filter *model.EventStreamFilter) (<-chan *supersense.Event, error) {
	pipe := make(chan *supersense.Event)

	if filter != nil {
		go r.mux.Register(pipe, ctx.Done(), filter.Sources...)
	} else {
		go r.mux.Register(pipe, ctx.Done())
	}

	return pipe, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
