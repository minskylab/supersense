package supersense

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Mux is a necessary struct to join different sources
type Mux struct {
	pipelines []chan Event
	sources []Source
	mu *sync.Mutex
}

// NewMux returns a new mux
func NewMux(ctx context.Context, sources ...Source) (*Mux, error) {
	channels := make([]chan Event, 0)
	m := &Mux{pipelines: channels, sources: sources, mu: &sync.Mutex{}}
	for _, s := range m.sources {
		go func(m *Mux, s Source) {
			for event := range s.Pipeline(ctx) {
				log.Warn(event.EmittedAt.Clock())
				for _, pipe := range m.pipelines {
					pipe <- event
				}
			}
		}(m, s)
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

// Pipeline returns the channel where arrive the all the events from the muxed sources
func (m *Mux) Pipeline() *chan Event {
	pipe := make(chan Event, 1)
	m.mu.Lock()
	m.pipelines = append(m.pipelines, pipe)
	m.mu.Unlock()
	return &m.pipelines[len(m.pipelines)-1]
}

func (m *Mux) Done(pipe *chan Event) {
	for i, p := range m.pipelines {
		if &p == pipe {
			m.mu.Lock()
			close(p)
			log.Error(m.pipelines)
			m.pipelines = append(m.pipelines[:i], m.pipelines[i+1:]...)
			log.Error(m.pipelines)
			m.mu.Lock()
		}
	}
}
