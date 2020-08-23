package boltdb

import (
	"time"

	"github.com/asdine/storm/v3"
	"github.com/minskylab/supersense"
	"github.com/pkg/errors"
)

func (s *Store) saveStateSnapshot(state *SnapshotSharedState) error {
	return s.db.Save(state)
}

func (s *Store) getStateSnapshot() (*SnapshotSharedState, error) {
	state := new(SnapshotSharedState)
	err := s.db.One("ID", s.mainID, state)
	if err != nil && err != storm.ErrNotFound {
		return nil, errors.WithStack(err)
	}
	return state, nil
}

func (s *Store) initialCheck() error {
	_, err := s.getStateSnapshot()
	if err != nil && err != storm.ErrNotFound {
		return errors.WithStack(err)
	}

	if err == storm.ErrNotFound {
		if err := s.saveStateSnapshot(&SnapshotSharedState{
			ID:        s.mainID,
			CreatedAt: time.Now(),
			UpdateAt:  time.Now(),
			TotalEvents: 0,
		}); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (s *Store) saveNewEvent(event Event) error {
	if err := s.db.Save(&event); err != nil {
		return errors.WithStack(err)
	}

	snapshot, err := s.getStateSnapshot()
	if err != nil {
		return errors.WithStack(err)
	}

	snapshot.TotalEvents += 1
	snapshot.UpdateAt = time.Now()

	return s.saveStateSnapshot(snapshot)
}

func (s *Store) getEvents(lasts int64) ([]*supersense.Event, error) {
	var board []Event

	if lasts < 1 || lasts > maxCurrentBoardBuffer {
		lasts = maxCurrentBoardBuffer
	}

	if err := s.db.AllByIndex("EmittedAt", &board, storm.Limit(int(lasts))); err != nil {
		return nil, errors.WithStack(err)
	}

	var finalEvents []*supersense.Event
	for _, e := range board {
		finalEvents = append(finalEvents, &e.Event)
	}

	return finalEvents, nil
}