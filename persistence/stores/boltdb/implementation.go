package boltdb

import (
	"time"

	"github.com/asdine/storm/v3"
	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/persistence"
	"github.com/pkg/errors"
)

const mainIDValue = "main_shared_state"

// SnapshotSharedState is a wrapper struct to make "storable" the SharedState
type SnapshotSharedState struct {
	ID                      string    `storm:"id"`
	CreatedAt               time.Time `json:"createdAt"`
	UpdateAt                time.Time `json:"updatedAt"`
	persistence.SharedState `storm:"inline"`
}

type Store struct {
	db        *storm.DB
	mainID    string
	maxBuffer int64
}

func NewStore(path string, maxBuffer int64) (*Store, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if maxBuffer <= 0 {
		maxBuffer = 100e3
	}

	return &Store{db: db, mainID: mainIDValue, maxBuffer: maxBuffer}, nil
}

func (s *Store) CurrentSharedState() (*persistence.SharedState, error) {
	state := new(SnapshotSharedState)
	if err := s.db.One("ID", s.mainID, state); err != nil {
		return nil, errors.WithStack(err)
	}

	return &state.SharedState, nil
}

func (s *Store) AddEventToSharedState(event supersense.Event) (*persistence.SharedState, error) {
	state := new(SnapshotSharedState)
	if err := s.db.One("ID", s.mainID, state); err != nil {
		return nil, errors.WithStack(err)
	}

	state.SharedState.Board = append(state.SharedState.Board, event)
	state.SharedState.LastUpdate = time.Now()

	if err := s.db.Save(state); err != nil {
		return nil, errors.WithStack(err)
	}

	return s.CurrentSharedState()
}

func (s *Store) Close() error {
	return s.db.Close()
}