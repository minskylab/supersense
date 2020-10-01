package sources

// "context"
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/minskylab/supersense"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const githubBaseURL = "https://api.github.com"
const rateProportionToLog = 1200

// Github is a source for three git repository events: Push, Fork, PullRequest
type Github struct {
	id               string
	sourceName       string
	channel          chan supersense.Event
	token            *string
	repos            []string
	eTags            map[string]string
	rateRemaining    map[string]string
	eventsDispatched []string // in memory state persistence
	httpClient       *http.Client
	baseURL          string
	mu               *sync.Mutex
	firstTime        bool
}

// NewGithub wraps all the needs for instance a new Github source
func NewGithub(token *string, repos []string) (*Github, error) {
	source := &Github{
		id:               uuid.NewV4().String(),
		sourceName:       "github",
		channel:          make(chan supersense.Event, 10),
		token:            token,
		repos:            repos,
		eTags:            map[string]string{},
		rateRemaining:    map[string]string{},
		eventsDispatched: []string{},
		baseURL:          githubBaseURL,
		firstTime:        true,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		mu: &sync.Mutex{},
	}
	return source, nil
}

// TODO: Pull Request: better titles
func (g *Github) pullEvents(owner, repo string, previousETag string, token *string) ([]*Event, *http.Response, error) {
	u := fmt.Sprintf("%s/repos/%v/%v/events", g.baseURL, owner, repo)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "go-github")
	req.Header.Set("If-None-Match", previousETag)

	if token != nil {
		req.Header.Set("Authorization", "token "+*token)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		log.Warnf("Github API response code: %d", resp.StatusCode)
	}

	var events []*Event

	decErr := json.NewDecoder(resp.Body).Decode(&events)
	if decErr != nil && decErr != io.EOF {
		return nil, nil, errors.WithStack(decErr)
	}

	return events, resp, nil
}

func (g *Github) alreadyDispatched(eventID string) bool {
	// g.mu.Lock()
	// defer g.mu.Unlock()

	for _, e := range g.eventsDispatched { // If the event has been dispatched
		if eventID == e {
			return true
		}
	}

	return false
}

func (g *GitHub) pullAndValidateEvents(parts []string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.mu.Lock()
	defer g.mu.Unlock()
	events, resp, err := g.pullEvents(parts[0], parts[1], g.eTags[repo], g.token)
	if err != nil {
		log.Errorf("%+v", err)
		return
	}

	// log.WithField("events", events).Debug("fetching " + repo)

	if resp == nil {
		log.Error("Invalid response from GitHub Events API.")
		return
	}

	etag := resp.Header.Get("ETag")
	etag = strings.TrimPrefix(etag, "W/")

	rateLimitRemaining := resp.Header.Get("X-Ratelimit-Remaining")
	// pollInterval := resp.Header.Get("X-Poll-Interval")

	g.eTags[repo] = etag

	if g.rateRemaining[repo] == "" {
		g.rateRemaining[repo] = rateLimitRemaining
	}

	rateRemaining, _ := strconv.Atoi(rateLimitRemaining)

	// for each multiple of rateProportionToLog, the logger prints a warning with the current rate limit
	if rateRemaining%rateProportionToLog == 0 || g.firstTime {
		log.WithFields(log.Fields{"etag": g.eTags[repo]}).Warn("Github API Rate Remaining: ", rateLimitRemaining)
	}

	if g.firstTime {
		g.firstTime = false
	}
}

func (g *Github) fetchRepo(repo string) {
	parts := strings.Split(repo, "/")

	events := g.pullAndValidateEvents(parts)

	for _, event := range events {
		if event.CreatedAt == nil || event.ID == nil || event.Type == nil {
			continue
		}

		if time.Now().Sub(*event.CreatedAt) > 100*time.Second {
			continue // No old events
		}

		if a := g.alreadyDispatched(*event.ID); a {
			continue
		}

		superEvent := supersense.Event{}
		superEvent.ID = *event.ID
		superEvent.CreatedAt = *event.CreatedAt
		superEvent.Actor = supersense.Person{}
		superEvent.SourceID = g.id
		superEvent.SourceName = g.sourceName

		// superEvent.ShareURL

		repo = strings.Trim(repo, "/")

		superEvent.Actor.Owner = g.sourceName
		repoLink := "https://github.com/" + repo
		superEvent.ShareURL = repoLink

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

			if superEvent.Actor.Username != nil {
				userRepoLink := "https://github.com/" + *superEvent.Actor.Username
				superEvent.Actor.ProfileURL = &userRepoLink
			}
		}

		switch payload.(type) {
		case *PushEvent:
			pushEvent := payload.(*PushEvent)
			for _, commit := range pushEvent.Commits {
				if commit.Message != nil {
					superEvent.Message = repo + ":\n" + *commit.Message
				}
			}

			if superEvent.Actor.Username != nil {
				superEvent.Title = "Push of " + *superEvent.Actor.Username
			} else {
				superEvent.Title = "Push"
			}

			if pushEvent.Pusher != nil {
				if pushEvent.Pusher.Login != nil {
					username := *pushEvent.Pusher.Login
					superEvent.Title += " of " + username
				}
			}
			superEvent.EventKind = "push"
		case *ForkEvent:
			forkEvent := payload.(*ForkEvent)

			if forkEvent.Forkee == nil {
				continue
			}

			forkeeRepo := ""

			if forkEvent.Forkee != nil {
				if forkEvent.Forkee != nil {
					forkeeRepo = *forkEvent.Forkee.FullName
				}
			}

			if forkeeRepo == "" && forkEvent.Forkee.Owner != nil {
				if forkEvent.Forkee.Owner.Login != nil {
					username := *forkEvent.Forkee.Owner.Login
					forkeeRepo += username
					superEvent.Title = "Fork of " + username
				}
			}

			if forkeeRepo == "" && forkEvent.Forkee.Name != nil {
				forkeeRepo += "/" + *forkEvent.Forkee.Name
			}

			superEvent.Message = repo + ":\n" + forkeeRepo
			superEvent.EventKind = "fork"

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

				message := title + ":\n" + body
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
		case *Issue:
			issueEvent := payload.(*Issue)

			var title, body, state, shareURL string
			if issueEvent.Title != nil {
				title = *issueEvent.Title
			}

			if issueEvent.Body != nil {
				body = *issueEvent.Body
			}

			if issueEvent.State != nil {
				state = *issueEvent.State
			}

			if issueEvent.HTMLURL != nil {
				shareURL = *issueEvent.HTMLURL
			}

			if shareURL == "" && issueEvent.URL != nil {
				shareURL = *issueEvent.URL
			}

			if issueEvent.User != nil {
				if issueEvent.User.Login != nil {
					ownerUsername := *issueEvent.User.Login
					superEvent.Title += " of " + ownerUsername
				}
			}

			superEvent.Title = repo + ":\n" + title
			superEvent.Message = body
			superEvent.EventKind = strings.Trim("new-issue-"+state, "- ")
			superEvent.ShareURL = shareURL

		case *IssuesEvent:
			issueEventWrap := payload.(*IssuesEvent)
			var action string
			if issueEventWrap.Action != nil {
				action = *issueEventWrap.Action
			}

			issueEvent := issueEventWrap.Issue
			if issueEvent == nil {
				continue
			}

			var title, shareURL string // body
			if issueEvent.Title != nil {
				title = *issueEvent.Title
			}

			// if issueEvent.Body != nil {
			// 	_ = *issueEvent.Body
			// }

			if issueEvent.HTMLURL != nil {
				shareURL = *issueEvent.HTMLURL
			}

			if shareURL == "" && issueEvent.URL != nil {
				shareURL = *issueEvent.URL
			}

			if issueEvent.User != nil {
				if issueEvent.User.Login != nil {
					ownerUsername := *issueEvent.User.Login
					superEvent.Title += " of " + ownerUsername
				}
			}

			superEvent.Title = "Issue Event"
			// superEvent.Message = body
			superEvent.Message = repo + ":\n" + title
			superEvent.EventKind = strings.Trim("issue-"+action, "- ")
			superEvent.ShareURL = shareURL
		default:
			// log.Debug(fmt.Sprintf("%T payload type not accepted in this stage of supersense", payload))
			continue
		}

		superEvent.EmittedAt = time.Now()
		g.mu.Lock()
		defer g.mu.Unlock()
		g.eventsDispatched = append(g.eventsDispatched, *event.ID)
		g.channel <- superEvent
	}
}

func (g *Github) loopFetchRepo(repo string) {
	for {
		g.fetchRepo(repo)
		time.Sleep(1000 * time.Millisecond)
	}
}

// Identify implements the Source interface
func (g *Github) Identify(nameOrID string) bool {
	return g.sourceName == nameOrID || g.id == nameOrID
}

// Run perform run initial procedure to spam the go-routine in charge to sniff the github events
func (g *Github) Run() error {
	for _, repo := range g.repos {
		go g.loopFetchRepo(repo)
	}

	return nil
}

// Pipeline returns the events channel
func (g *Github) Pipeline() <-chan supersense.Event {
	return g.channel
}

// Dispose close all streams and flows with the source
func (g *Github) Dispose() {
	close(g.channel)
}
