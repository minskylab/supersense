package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/graph/generated"
	"github.com/minskylab/supersense/graph/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*model.AuthResponse, error) {
	token, err := r.db.LoginWithUserPassword(username, password)
	if err != nil {
		return nil, errors.WithMessage(err, "invalid token check more options")
	}
	return &model.AuthResponse{
		JwtToken: token,
		ExpirateAt: time.Now(),
	}, nil
}

func (r *mutationResolver) Broadcast(ctx context.Context, auth string, draft model.EventDraft) (string, error) {
	user := ctx.Value("user")
	logrus.Info(spew.Sdump(user))
	return "", nil
}

func (r *queryResolver) Event(ctx context.Context, id string) (*supersense.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *subscriptionResolver) Events(ctx context.Context) (<-chan *supersense.Event, error) {
	pipe := make(chan *supersense.Event, 1)

	go func() {
		for event := range r.mux.Events() {
			pipe <- &event
		}
	}()

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
