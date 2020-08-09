package sources

// "context"
import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	// "github.com/google/go-github/v32/github"

	"github.com/google/go-github/github"
	"github.com/minskylab/supersense"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// Github is a source for three git repository events: Push, Fork, PullRequest
type Github struct {
	id               string
	name             string
	channel          chan supersense.Event
	token            *string
	repos            []string
	etags            map[string]string
	rateRemainings   map[string]string
	eventsDispatched []string
}

// NewGithub wraps all the needs for instace a new Github source
func NewGithub(repos []string, token *string) (*Github, error) {
	source := &Github{
		id:               uuid.NewV4().String(),
		name:             "github",
		channel:          make(chan supersense.Event, 1),
		token:            token,
		repos:            repos,
		etags:            map[string]string{},
		rateRemainings:   map[string]string{},
		eventsDispatched: []string{},
	}
	return source, nil
}

// Run perform run initial procedure to spam the gorutine in charge to sniff the github events
func (g *Github) Run(ctx context.Context) error {
	var httpClient *http.Client = nil

	if g.token != nil {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: *g.token},
		)
		httpClient = oauth2.NewClient(ctx, ts)
	}

	client := github.NewClient(httpClient)

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
				// pollInterval := resp.Header.Get("X-Poll-Interval")

				if g.etags[repo] == "" {
					g.etags[repo] = etag
				}

				if g.rateRemainings[repo] == "" {
					g.rateRemainings[repo] = rateLimitRemaining
				}

				// log.Info("etag: ", etag)
				// log.Info("pollInterval: ", pollInterval)
				// log.Info("rateLimitRemaining: ", rateLimitRemaining)
				rateRemaining, _ := strconv.Atoi(rateLimitRemaining)

				if rateRemaining%1200 == 0 {
					log.Info("Github API Rate Remaining: ", rateLimitRemaining)
				}

				for _, event := range events {

					if time.Now().Sub(event.GetCreatedAt()) > 6*time.Second {
						continue // No old events
					}

					eventID := event.GetID() // If the event has been dispatched
					for _, e := range g.eventsDispatched {
						if eventID == e {
							continue
						}
					}

					log.Info("Github event type: " + event.GetType())

					superEvent := supersense.Event{}
					superEvent.ID = event.GetID()
					superEvent.CreatedAt = event.GetCreatedAt()
					superEvent.Actor = supersense.Person{}
					superEvent.SourceID = g.id
					superEvent.SourceName = g.name

					// superEvent.ShareURL

					superEvent.Actor.Owner = g.name
					repoLink := "https://github.com/" + repo
					superEvent.Actor.ProfileURL = &repoLink

					if event == nil {
						continue
					}

					// event.Payload()
					payload, err := event.ParsePayload()
					if err != nil {
						log.Warn(errors.WithStack(err))
					}
					// log.Info("")
					// log.Info("[TYPE] ", event.GetType())

					// superEvent.EventKind = event.GetType()

					if event.GetActor() != nil {
						// log.Info("[ACTOR] [NAME]", event.GetActor().GetName())
						superEvent.Actor.Name = event.GetActor().GetName()
						// log.Info("[ACTOR] [EMAIL]", event.GetActor().GetEmail())
						superEvent.Actor.Email = event.GetActor().Email
						// log.Info("[ACTOR] [USERNAME]", event.GetActor().GetLogin())
						superEvent.Actor.Username = event.GetActor().Login
						// log.Info("[ACTOR] [AVATAR]", event.GetActor().GetAvatarURL())
						superEvent.Actor.Photo = event.GetActor().GetAvatarURL()
					}

					switch payload.(type) {
					case *github.PushEvent:
						pushEvent := payload.(*github.PushEvent)
						for _, commit := range pushEvent.Commits {
							// log.Info("[PUSH] [COMMIT] ", commit.GetMessage())
							superEvent.Message = commit.GetMessage()
						}
						superEvent.Title = "Push"
						if pushEvent.GetPusher() != nil {
							username := pushEvent.GetPusher().GetLogin()
							superEvent.Title += " of " + username
						}
						superEvent.EventKind = "push"
					case *github.ForkEvent:
						forkEvent := payload.(*github.ForkEvent)

						if forkEvent.GetForkee() != nil {
							forkeeRepo := ""
							if forkEvent.GetForkee().GetOwner() != nil {
								// log.Info("[FORK] [FORKEE_OWNER_USERNAME] ", forkEvent.GetForkee().GetOwner().GetLogin())
								username := forkEvent.GetForkee().GetOwner().GetLogin()
								forkeeRepo += username
								superEvent.Title = "Fork of " + username
							}
							// log.Info("[FORK] [FORKEE_NAME] ", forkEvent.GetForkee().GetName())
							forkeeRepo += forkEvent.GetForkee().GetName()
							superEvent.Message = forkeeRepo
							superEvent.EventKind = "fork"
						}
					case *github.PullRequestEvent:
						pullRequestEvent := payload.(*github.PullRequestEvent)
						pullRequest := pullRequestEvent.GetPullRequest()
						if pullRequest != nil {
							// log.Info("[FORK] [PULLREQUEST_TITLE] ", pullRequest.GetTitle())
							title := pullRequest.GetTitle()
							// log.Info("[FORK] [PULLREQUEST_BODY] ", pullRequest.GetBody())
							body := pullRequest.GetBody()
							// log.Info("[FORK] [PULLREQUEST_STATE] ", pullRequest.GetState())
							state := pullRequest.GetState()

							message := title + "\n" + body
							superEvent.Message = message

							superEvent.EventKind = strings.Trim("pull-request-"+state, "- ")
							if pullRequest.GetUser() != nil {
								// log.Info("[FORK] [PULLREQUEST_OWNER_NAME] ", pullRequest.GetUser().GetName())
								// log.Info("[FORK] [PULLREQUEST_OWNER_USERNAME] ", pullRequest.GetUser().GetLogin())
								ownerUsername := pullRequest.GetUser().GetLogin()
								// log.Info("[FORK] [PULLREQUEST_OWNER_EMAIL]", pullRequest.GetUser().GetEmail())

								superEvent.Title = "Pull Request of " + ownerUsername
							}

						}

					default:
						log.Error(fmt.Sprintf("%T", payload), " payload type not accepted")
					}

					superEvent.EmmitedAt = time.Now()
					g.eventsDispatched = append(g.eventsDispatched, eventID)
					g.channel <- superEvent
				}
			}
			time.Sleep(3 * time.Second)
		}
	}()

	return nil
}

// Events returns the events channel
func (g *Github) Events(ctx context.Context) *chan supersense.Event {
	return &g.channel
}
