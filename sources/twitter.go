package sources

import (
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
	token := oauth1.NewToken(props.AccessToken, props.AccessToken)

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
func (s *Twitter) Run() error {
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		*s.events <- supersense.Event{
			ID:        tweet.IDStr,
			EmmitedAt: time.Now(),
			Message:   tweet.Text,
			SourceID:  s.id,
			Title:     fmt.Sprintf("Tweet from %s", tweet.User.Name),
			Person: supersense.Person{
				Name:        tweet.User.Name,
				Email:       &tweet.User.Email,
				Photo:       tweet.User.ProfileBackgroundImageURL,
				SourceOwner: s.sourceName,
			},
		}
		fmt.Println(tweet.Text)
	}
	demux.DM = func(dm *twitter.DirectMessage) {
		log.Infof("(DM) %s | %s", dm.SenderID, dm.Text)
	}
	demux.Event = func(event *twitter.Event) {
		log.Infof("(EV) %s | %s", event.Source.ID, event.Event)
	}

	stream, err := s.client.Streams.Filter(&twitter.StreamFilterParams{
		Track:         s.queryToTrack,
		StallWarnings: twitter.Bool(true),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	go demux.HandleChan(stream.Messages)

	return nil
}

// Events return a channel from where come in the events
func (s *Twitter) Events() *chan supersense.Event {
	return s.events
}
