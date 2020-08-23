package persistence

import "github.com/minskylab/supersense"

// Store represents an interface where you can preserve your global state
type Store interface {
	CurrentSharedState(lasts int64) (*SharedState, error)
	AddEventToSharedState(event *supersense.Event) error
	Close() error
}
