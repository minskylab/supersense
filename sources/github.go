package sources

// "context"
import (
	"context"

	// "github.com/google/go-github/v32/github"

	"github.com/minskylab/supersense"
)

type Github struct {
}

func NewGithub() {
	// TODO: Search about Event API and long pooling
	// client := github.NewClient(nil)
	// client.Activity.ListRepositoryEvents(ctx context.Context, owner string, repo string, opts *github.ListOptions)

}

func (g *Github) Run(ctx context.Context) error {
	return nil
}

func (g *Github) Events(ctx context.Context) *chan supersense.Event {

	return nil
}
