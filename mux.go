package supersense

import (
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Mux is a necessary struct to join different sources
type Mux struct {
	channel chan Event
	sources []Source
}

// NewMux returns a new mux
func NewMux(ctx context.Context, sources ...Source) (*Mux, error) {
	generalChannel := make(chan Event, 10)
	m := &Mux{channel: generalChannel, sources: sources}
	for _, s := range m.sources {
		go func(s Source) {
			for event := range *s.Events(ctx) {
				log.Warn(event.EmittedAt.Clock())
				m.channel <- event
			}
		}(s)
	}
	return m, nil
}

// RunAllSources run all the sources at the same time
func (m *Mux) RunAllSources(ctx context.Context) error {
	for _, s := range m.sources {
		if err := s.Run(ctx); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Events returns the channel where arrive the all the events from the muxed sources
func (m *Mux) Events() chan Event {
	return m.channel
}
