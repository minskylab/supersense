package sources

import (
	"fmt"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/minskylab/supersense"
	log "github.com/sirupsen/logrus"
)

// Twitter represents a twitter stream source
type Twitter struct {
	id           string
	sourceName   string
	queryToTrack []string
	client       *twitter.Client
	events       chan supersense.Event
}

// TwitterClientProps wraps minimal information for a twitter client creation
type TwitterClientProps struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
	QueryToTrack   []string
}

// NewTwitter returns a new Twitter Source
func NewTwitter(props TwitterClientProps) (*Twitter, error) {
	config := oauth1.NewConfig(props.ConsumerKey, props.ConsumerSecret)
	token := oauth1.NewToken(props.AccessToken, props.AccessSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	httpClient.Timeout = 60 * time.Second
	client := twitter.NewClient(httpClient)
	eventsChan := make(chan supersense.Event, 1)

	return &Twitter{
		id:           uuid.NewV4().String(),
		sourceName:   "twitter",
		queryToTrack: props.QueryToTrack,
		client:       client,
		events:       eventsChan,
	}, nil
}

// Identify implements the Source interface
func (s *Twitter) Identify(nameOrID string) bool {
	return s.sourceName == nameOrID || s.id == nameOrID
}

// optimize later
func urlAlreadyExists(urls []supersense.URLEntity, url string) bool {
	for _, u := range urls {
		if u.URL == url {
			return true
		}
	}
	return false
}

// optimize later
func mediaAlreadyExists(medias []supersense.MediaEntity, mediaURL string) bool {
	for _, m := range medias {
		if m.URL == mediaURL {
			return true
		}
	}
	return false
}

// optimize later
func tagAlreadyExists(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

// Run bootstrap the necessary actions to demux and listen new tweets
// and implements the Source interface of supersense
func (s *Twitter) Run() error {
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		entities := supersense.Entities{
			Urls:  []supersense.URLEntity{},
			Media: []supersense.MediaEntity{},
			Tags:  []string{},
		}

		if tweet.ExtendedTweet != nil && tweet.ExtendedTweet.Entities != nil {
			ents := tweet.ExtendedTweet
			for _, url := range ents.Entities.Urls {
				entities.Urls = append(entities.Urls, supersense.URLEntity{
					DisplayURL: url.DisplayURL,
					URL:        url.URL,
				})
			}

			for _, media := range ents.Entities.Media {
				entities.Media = append(entities.Media, supersense.MediaEntity{
					Type: media.Type,
					URL:  media.MediaURLHttps,
				})
			}

			for _, tag := range ents.Entities.Hashtags {
				entities.Tags = append(entities.Tags, tag.Text)
			}
		}

		// more content to original tweet
		// but with lower priority
		if tweet.Retweeted {
			if tweet.RetweetedStatus != nil && tweet.RetweetedStatus.Entities != nil {
				ents := tweet.RetweetedStatus.Entities
				for _, url := range ents.Urls {
					if !urlAlreadyExists(entities.Urls, url.URL) {
						entities.Urls = append(entities.Urls, supersense.URLEntity{
							DisplayURL: url.DisplayURL,
							URL:        url.URL,
						})
					}
				}

				for _, media := range ents.Media {
					if !mediaAlreadyExists(entities.Media, media.MediaURLHttps) {
						entities.Media = append(entities.Media, supersense.MediaEntity{
							Type: media.Type,
							URL:  media.MediaURLHttps,
						})
					}
				}

				for _, tag := range ents.Hashtags {
					if !tagAlreadyExists(entities.Tags, tag.Text) {
						entities.Tags = append(entities.Tags, tag.Text)
					}
				}
			}
		}

		message := tweet.Text
		if tweet.ExtendedTweet != nil {
			message = tweet.ExtendedTweet.FullText
		}

		createdAt, _ := time.Parse(time.RubyDate, tweet.CreatedAt)
		person := supersense.Person{}

		if tweet.User != nil {
			person.Name = tweet.User.Name
			person.Photo = tweet.User.ProfileImageURLHttps
			person.Owner = s.sourceName
			person.Email = &tweet.User.Email
			person.ProfileURL = &tweet.User.URL
			person.Username = &tweet.User.ScreenName
		}

		eventTweetKind := "tweet"
		if strings.HasPrefix(strings.TrimSpace(message), "RT") {
			eventTweetKind = "retweet"
		}

		s.events <- supersense.Event{
			ID:         tweet.IDStr,
			CreatedAt:  createdAt,
			EmittedAt:  time.Now(),
			Message:    message,
			SourceID:   s.id,
			SourceName: s.sourceName,
			EventKind:  eventTweetKind,
			Title:      fmt.Sprintf("%s of %s", eventTweetKind, tweet.User.Name),
			Entities:   entities,
			// ShareURL:   tweet.Source,
			ShareURL: fmt.Sprintf("https://twitter.com/i/web/status/%s", tweet.IDStr),
			Actor:    person,
		}

	}

	demux.DM = func(dm *twitter.DirectMessage) {
		log.Infof("(DM) %s | %s", dm.SenderID, dm.Text)
	}

	demux.Event = func(event *twitter.Event) {
		log.Warnf("(EV) %s | %s", event.Source.ID, event.Event)
	}

	stream, err := s.client.Streams.Filter(&twitter.StreamFilterParams{
		StallWarnings: twitter.Bool(true),
		Track:         s.queryToTrack,
	})

	if err != nil {
		return errors.WithStack(err)
	}

	go func() {
		demux.HandleChan(stream.Messages)
	}()

	return nil
}

// Pipeline return a channel from where come in the events
func (s *Twitter) Pipeline() <-chan supersense.Event {
	return s.events
}

// Dispose return a channel from where come in the events
func (s *Twitter) Dispose() {
	return
}
