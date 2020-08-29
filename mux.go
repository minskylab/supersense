package supersense

import (
	"sync"

	"github.com/pkg/errors"
)

// Mux is a necessary struct to join different sources
type Mux struct {
	pipelines []chan *Event            // fan-out pipelines
	filters   map[chan *Event][]string // sources filter
	sources   []Source
	running   map[Source]bool
	mu        *sync.Mutex
}

// NewMux returns a new mux to use as a mani pipeline for all your event sources
func NewMux(sources ...Source) (*Mux, error) {
	channels := make([]chan *Event, 0)
	m := &Mux{
		mu:        &sync.Mutex{},
		sources:   sources,
		pipelines: channels,
		running:   map[Source]bool{},
		filters:   map[chan *Event][]string{},
	}
	return m, nil
}

func (m *Mux) setRunningSource(s Source, isRunning bool) {
	m.mu.Lock()
	m.running[s] = isRunning
	m.mu.Unlock()
}

func (m *Mux) sourceListener(s Source) {
	m.setRunningSource(s, true)
	for event := range s.Pipeline() {
		m.mu.Lock()
		for _, pipe := range m.pipelines {
			filters, filtered := m.filters[pipe]
			if filtered && len(filters) > 0 {
				for _, filter := range filters {
					if filter == event.SourceName {
						pipe <- &event
					}
				}
			} else {
				pipe <- &event
			}
		}
		m.mu.Unlock()
	}
	m.setRunningSource(s, false)
}

func (m *Mux) addNewSource(s Source) {
	m.sources = append(m.sources, s)
	go m.sourceListener(s)
}

// AddNewSource exports this function to public
func (m *Mux) AddNewSource(s Source) {
	m.addNewSource(s)
}

// RunAllSources run all the sources at the same time
func (m *Mux) RunAllSources() error {
	for _, s := range m.sources {
		go m.sourceListener(s)
	}

	for _, s := range m.sources {
		if err := s.Run(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (m *Mux) addPipeline(pipeline chan *Event, filteredSources ...string) {
	m.mu.Lock()
	m.pipelines = append(m.pipelines, pipeline)
	if len(filteredSources) > 0 {
		m.filters[pipeline] = filteredSources
	}
	m.mu.Unlock()
}

func (m *Mux) deleteAndClosePipeline(pipeline chan *Event) {
	for i, p := range m.pipelines {
		if p == pipeline {
			m.mu.Lock()
			close(p)
			m.pipelines = append(m.pipelines[:i], m.pipelines[i+1:]...)
			delete(m.filters, pipeline)
			m.mu.Unlock()
		}
	}
}

// Register attach a new channel to the pipes list.
func (m *Mux) Register(pipeline chan *Event, done <-chan struct{}, filteredSources ...string) {
	m.addPipeline(pipeline, filteredSources...)
	<-done
	m.deleteAndClosePipeline(pipeline)
}
