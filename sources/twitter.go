package sources

import (
	"context"
	"fmt"
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
	events       *chan supersense.Event
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

	client := twitter.NewClient(httpClient)
	eventsChan := make(chan supersense.Event, 1)
	return &Twitter{
		id:           uuid.NewV4().String(),
		sourceName:   "twitter",
		queryToTrack: props.QueryToTrack,
		client:       client,
		events:       &eventsChan,
	}, nil
}

// Run bootstrap the necessary actions to demux and listen new tweets
// and implements the Source interface of supersense
func (s *Twitter) Run(ctx context.Context) error {
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		entities := supersense.Entities{
			Urls:  []supersense.URLEntity{},
			Media: []supersense.MediaEntity{},
			Tags:  []string{},
		}
		if tweet.ExtendedTweet != nil {
			if tweet.ExtendedTweet.Entities != nil {
				for _, url := range tweet.ExtendedTweet.Entities.Urls {
					entities.Urls = append(entities.Urls, supersense.URLEntity{
						DisplayURL: url.DisplayURL,
						URL:        url.URL,
					})
				}

				for _, media := range tweet.ExtendedTweet.Entities.Media {
					entities.Media = append(entities.Media, supersense.MediaEntity{
						Type: media.Type,
						URL:  media.MediaURLHttps,
					})
				}

				for _, tag := range tweet.ExtendedTweet.Entities.Hashtags {
					entities.Tags = append(entities.Tags, tag.Text)
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

		*s.events <- supersense.Event{
			ID:         tweet.IDStr,
			CreatedAt:  createdAt,
			EmmitedAt:  time.Now(),
			Message:    message,
			SourceID:   s.id,
			SourceName: s.sourceName,
			EventKind:  "tweet",
			Title:      fmt.Sprintf("Tweet of %s", tweet.User.Name),
			Entities:   entities,
			ShareURL:   tweet.Source,
			Person:     person,
		}

	}
	demux.DM = func(dm *twitter.DirectMessage) {
		log.Infof("(DM) %s | %s", dm.SenderID, dm.Text)
	}
	demux.Event = func(event *twitter.Event) {
		log.Warnf("(EV) %s | %s", event.Source.ID, event.Event)
	}

	stream, err := s.client.Streams.Filter(&twitter.StreamFilterParams{
		Track:         s.queryToTrack,
		StallWarnings: twitter.Bool(true),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	go func() {
		demux.HandleChan(stream.Messages)
	}()

	return nil
}

// Events return a channel from where come in the events
func (s *Twitter) Events(ctx context.Context) *chan supersense.Event {
	return s.events
}
