package boltdb

import (
	"time"

	"github.com/asdine/storm/v3"
	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/persistence"
	"github.com/pkg/errors"
)


// SnapshotSharedState is a wrapper struct to make "storable" the SharedState
type SnapshotSharedState struct {
	ID                      string    `storm:"id"`
	CreatedAt               time.Time `json:"createdAt"`
	UpdateAt                time.Time `json:"updatedAt"`
	TotalEvents int64 `json:"total_events"`
}

type Store struct {
	db        *storm.DB
	mainID    string
	maxBuffer int64
}

func (s *Store) CurrentSharedState(lasts int64) (*persistence.SharedState, error) {
	snapshot, err := s.getStateSnapshot()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	events, err := s.getEvents(lasts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	state := &persistence.SharedState{
		Board:     events,
		LastUpdate: snapshot.UpdateAt,
	}

	return state, nil
}

func (s *Store) AddEventToSharedState(event *supersense.Event) error {
	return s.saveNewEvent(Event{Event: *event, ID: event.ID, EmittedAt: event.EmittedAt})
}

func (s *Store) Close() error {
	return s.db.Close()
}
