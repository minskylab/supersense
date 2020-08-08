package sources

// "context"
import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	// "github.com/google/go-github/v32/github"

	"github.com/google/go-github/github"
	"github.com/minskylab/supersense"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Github struct {
	repos          []string
	etags          map[string]string
	rateRemainings map[string]string
}

func NewGithub(repos []string) (*Github, error) {
	source := &Github{
		repos:          repos,
		etags:          map[string]string{},
		rateRemainings: map[string]string{},
	}
	return source, nil
}

func (g *Github) Run(ctx context.Context) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	go func() {
		for {
			for _, repo := range g.repos {
				// repo := "minskylab/supersense"
				parts := strings.Split(repo, "/")

				events, resp, err := client.Activity.ListRepositoryEvents(ctx, parts[0], parts[1],
					&github.ListOptions{
						PerPage: 250,
					})
				if err != nil {
					panic(err)
				}

				etag := resp.Header.Get("ETag")
				rateLimitRemaining := resp.Header.Get("X-Ratelimit-Remaining")
				pollInterval := resp.Header.Get("X-Poll-Interval")

				if g.etags[repo] == "" {
					g.etags[repo] = etag
				}

				if g.rateRemainings[repo] == "" {
					g.rateRemainings[repo] = rateLimitRemaining
				}

				log.Info("etag: ", etag)
				log.Info("pollInterval: ", pollInterval)
				log.Info("rateLimitRemaining: ", rateLimitRemaining)

				for _, event := range events {
					if event.CreatedAt != nil {
						if time.Now().Sub(*event.CreatedAt) > 2*time.Second {
							continue // No old events
						}
					}

					if event != nil {
						// event.Payload()
						payload, err := event.ParsePayload()
						if err != nil {
							log.Warn(errors.WithStack(err))
						}
						log.Info("")
						log.Info("[TYPE] ", event.GetType())

						if event.GetActor() != nil {
							log.Info("[ACTOR] [NAME]", event.GetActor().GetName())
							log.Info("[ACTOR] [EMAIL]", event.GetActor().GetEmail())
							log.Info("[ACTOR] [USERNAME]", event.GetActor().GetLogin())
							log.Info("[ACTOR] [AVATAR]", event.GetActor().GetAvatarURL())
						}

						switch payload.(type) {
						case *github.PushEvent:
							pushEvent := payload.(*github.PushEvent)
							for _, commit := range pushEvent.Commits {
								log.Info("[PUSH] [COMMIT] ", commit.GetMessage())
							}

							if pushEvent.GetPusher() != nil {
								log.Info("[PUSH] [PUSHER_NAME] ", pushEvent.GetPusher().GetName())
								log.Info("[PUSH] [PUSHER_EMAIL] ", pushEvent.GetPusher().GetEmail())
								log.Info("[PUSH] [PUSHER_USERNAME] ", pushEvent.GetPusher().GetLogin())
							}
						case *github.ForkEvent:
							forkEvent := payload.(*github.ForkEvent)

							if forkEvent.GetForkee() != nil {
								log.Info("[FORK] [FORKEE_NAME] ", forkEvent.GetForkee().GetName())
								if forkEvent.GetForkee().GetOwner() != nil {
									log.Info("[FORK] [FORKEE_OWNER_NAME] ", forkEvent.GetForkee().GetOwner().GetName())
									log.Info("[FORK] [FORKEE_OWNER_USERNAME] ", forkEvent.GetForkee().GetOwner().GetLogin())
									log.Info("[FORK] [FORKEE_OWNER_EMAIL] ", forkEvent.GetForkee().GetOwner().GetEmail())
								}
								if forkEvent.GetSender() != nil {
									log.Info("[FORK] [FORKEE_SENDER_NAME] ", forkEvent.GetSender().GetName())
									log.Info("[FORK] [FORKEE_SENDER_USERNAME] ", forkEvent.GetSender().GetLogin())
									log.Info("[FORK] [FORKEE_SENDER_EMAIL] ", forkEvent.GetSender().GetEmail())
								}
							}
						case *github.PullRequestEvent:
							pullRequestEvent := payload.(*github.PullRequestEvent)
							pullRequest := pullRequestEvent.GetPullRequest()
							if pullRequest != nil {
								log.Info("[FORK] [PULLREQUEST_TITLE] ", pullRequest.GetTitle())
								log.Info("[FORK] [PULLREQUEST_BODY] ", pullRequest.GetBody())
								log.Info("[FORK] [PULLREQUEST_STATE] ", pullRequest.GetState())

								if pullRequest.GetUser() != nil {
									log.Info("[FORK] [PULLREQUEST_OWNER_NAME] ", pullRequest.GetUser().GetName())
									log.Info("[FORK] [PULLREQUEST_OWNER_USERNAME] ", pullRequest.GetUser().GetLogin())
									log.Info("[FORK] [PULLREQUEST_OWNER_EMAIL]", pullRequest.GetUser().GetEmail())
								}

							}

						default:
							log.Error(fmt.Sprintf("%T", payload), " payload type not accepted")
						}

					}
				}
			}
			time.Sleep(3 * time.Second)
		}
	}()

	return nil
}

func (g *Github) Events(ctx context.Context) *chan supersense.Event {

	return nil
}
