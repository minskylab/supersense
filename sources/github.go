package sources

// "context"
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/minskylab/supersense"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const githubBaseURL = "https://api.github.com"

// Github is a source for three git repository events: Push, Fork, PullRequest
type Github struct {
	id               string
	name             string
	channel          chan supersense.Event
	token            *string
	repos            []string
	eTags            map[string]string
	rateRemaining    map[string]string
	eventsDispatched []string // in memory state persistence
	httpClient       *http.Client
	baseURL          string
	mu *sync.Mutex
}

// NewGithub wraps all the needs for instance a new Github source
func NewGithub(token *string, repos []string) (*Github, error) {
	source := &Github{
		id:               uuid.NewV4().String(),
		name:             "github",
		channel:          make(chan supersense.Event, 10),
		token:            token,
		repos:            repos,
		eTags:            map[string]string{},
		rateRemaining:    map[string]string{},
		eventsDispatched: []string{},
		baseURL: githubBaseURL,
		httpClient: &http.Client{
			Timeout:       60*time.Second,
		},
		mu: &sync.Mutex{},
	}
	return source, nil
}

// TODO: Pull Request: better titles

func (g *Github) pullEvents(owner, repo string, previousETag string, token *string) ([]*Event, *http.Response, error){
	u := fmt.Sprintf("%s/repos/%v/%v/events", g.baseURL, owner, repo)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	req.Header.Set("Accept",  "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "go-github")
	req.Header.Set("If-None-Match", previousETag)

	log.Info("If-None-Match: ", previousETag)

	if token != nil {
		req.Header.Set("Authorization", "token " + *token)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, resp, errors.WithStack(err)
	}

	var events []*Event

	decErr := json.NewDecoder(resp.Body).Decode(&events)
	if decErr != nil && decErr != io.EOF {
		return nil, nil, errors.WithStack(err)
	}

	return events, resp, nil
}

func (g *Github) fetchRepo(repo string) {
	parts := strings.Split(repo, "/")
	events, resp, err := g.pullEvents(parts[0], parts[1], g.eTags[repo], g.token)
	if err != nil {
		log.Errorf("%+v", err)
		return
	}

	etag := resp.Header.Get("ETag")
	etag = strings.TrimPrefix(etag, "W/")

	rateLimitRemaining := resp.Header.Get("X-Ratelimit-Remaining")
	// pollInterval := resp.Header.Get("X-Poll-Interval")

	g.mu.Lock()
	if g.eTags[repo] == "" {
		g.eTags[repo] = etag
	}

	if g.rateRemaining[repo] == "" {
		g.rateRemaining[repo] = rateLimitRemaining
	}

	// rateRemaining, _ := strconv.Atoi(rateLimitRemaining)

	// if rateRemaining%1200 == 0 {
	log.WithFields(log.Fields{"repo": repo, "etag": g.eTags[repo]}).Warn("Github API Rate Remaining: ", rateLimitRemaining)
	// }

	g.mu.Unlock()

	for _, event := range events {
		if event.CreatedAt == nil || event.ID == nil || event.Type == nil {
			continue
		}

		if time.Now().Sub(*event.CreatedAt) > 6*time.Second {
			continue // No old events
		}

		for _, e := range g.eventsDispatched { // If the event has been dispatched
			if *event.ID == e {
				continue
			}
		}

		log.Info("Github event type: " + *event.Type)

		superEvent := supersense.Event{}
		superEvent.ID = *event.ID
		superEvent.CreatedAt = *event.CreatedAt
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

		if event.Actor != nil {
			if event.Actor.Name != nil {
				superEvent.Actor.Name = *event.Actor.Name
			}
			if event.Actor.AvatarURL != nil {
				superEvent.Actor.Photo = *event.Actor.AvatarURL
			}
			superEvent.Actor.Email = event.Actor.Email
			superEvent.Actor.Username = event.Actor.Login
		}

		switch payload.(type) {
		case *PushEvent:
			pushEvent := payload.(*PushEvent)
			for _, commit := range pushEvent.Commits {
				if commit.Message != nil {
					superEvent.Message = *commit.Message
				}
			}
			superEvent.Title = "Push"
			if pushEvent.Pusher != nil {
				if pushEvent.Pusher.Login != nil {
					username := *pushEvent.Pusher.Login
					superEvent.Title += " of " + username
				}
			}
			superEvent.EventKind = "push"
		case *ForkEvent:
			forkEvent := payload.(*ForkEvent)

			if forkEvent.Forkee != nil {
				forkeeRepo := ""
				if forkEvent.Forkee.Owner != nil {
					if forkEvent.Forkee.Owner.Login != nil {
						username := *forkEvent.Forkee.Owner.Login
						forkeeRepo += username
						superEvent.Title = "Fork of " + username
					}
				}

				if forkEvent.Forkee.Name != nil {
					forkeeRepo += "/" + *forkEvent.Forkee.Name
				}

				superEvent.Message = forkeeRepo
				superEvent.EventKind = "fork"
			}
		case *PullRequestEvent:
			pullRequestEvent := payload.(*PullRequestEvent)
			pullRequest := pullRequestEvent.PullRequest
			if pullRequest != nil {
				var title, body, state string
				if pullRequest.Title != nil {
					title = *pullRequest.Title
				}

				if pullRequest.Body != nil {
					body = *pullRequest.Body
				}

				if pullRequest.State != nil {
					state = *pullRequest.State
				}

				message := title + "\n" + body
				superEvent.Message = message

				superEvent.EventKind = strings.Trim("pull-request-"+state, "- ")

				if pullRequest.User != nil {
					superEvent.Title = "Pull Request"
					if pullRequest.User.Login != nil {
						ownerUsername := *pullRequest.User.Login
						superEvent.Title += " of " + ownerUsername
					}
				}

			}

		default:
			log.Warn(fmt.Sprintf("%T", payload), " payload type not accepted in this stage of supersense")
		}

		superEvent.EmittedAt = time.Now()
		g.mu.Lock()
		g.eventsDispatched = append(g.eventsDispatched, *event.ID)
		g.channel <- superEvent
		g.mu.Unlock()
	}
}

func (g *Github) loopFetchRepo(repo string) {
	for {
		g.fetchRepo(repo)
		time.Sleep(1 * time.Second)
	}
}

// Run perform run initial procedure to spam the go-routine in charge to sniff the github events
func (g *Github) Run(ctx context.Context) error {
	for _, repo := range g.repos {
		go g.loopFetchRepo(repo)
	}

	return nil
}

// Pipeline returns the events channel
func (g *Github) Pipeline(ctx context.Context) <-chan supersense.Event {
	return g.channel
}

// Dispose close all streams and flows with the source
func (g *Github) Dispose(ctx context.Context) {
	return
}
