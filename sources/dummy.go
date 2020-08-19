package sources

import (
	"context"
	"time"

	"github.com/minskylab/supersense"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// Dummy is a minimal source implementation,
// it's util when you need to test supersense
type Dummy struct {
	id         string
	sourceName string
	clock      *time.Ticker
	message    string
	events     *chan supersense.Event
}

// NewDummy creates and init a new dummy source
func NewDummy(period time.Duration, message string) (*Dummy, error) {
	eventsChan := make(chan supersense.Event, 10)
	source := &Dummy{
		id:         uuid.NewV4().String(),
		sourceName: "dummy",
		clock:      time.NewTicker(period),
		events:     &eventsChan,
		message:    message,
	}
	return source, nil
}

// Run starts the recurrent message issuer
func (s *Dummy) Run(ctx context.Context) error {
	if s.events == nil {
		return errors.New("invalid Source, it not have an events channel")
	}
	username := "jhondoe"

	go func() {
		for {
			event := <-s.clock.C
			*s.events <- supersense.Event{
				ID:         uuid.NewV4().String(),
				Message:    s.message,
				EmittedAt:  event,
				SourceID:   s.id,
				SourceName: s.sourceName,
				EventKind:  "dummy",
				ShareURL:   "https://example.com",
				Actor: supersense.Person{
					Name:     "John Doe",
					Photo:    "https://pic.jpeg",
					Owner:    s.sourceName,
					Username: &username,
				},
			}
		}
	}()

	return nil
}

// Events implements the supersense.Source interface
func (s *Dummy) Events(ctx context.Context) *chan supersense.Event {
	return s.events
}

// Dispose close all streams and flows with the source
func (s *Dummy) Dispose(ctx context.Context) {
	return
}