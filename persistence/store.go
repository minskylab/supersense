package persistence

import "github.com/minskylab/supersense"

// Store represents an interface where you can preserve your global state
type Store interface {
	CurrentSharedState() (*SharedState, error)
	AddEventToSharedState(event *supersense.Event) (*SharedState, error)
	Close() error
}
