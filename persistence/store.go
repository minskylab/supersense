package persistence

import "github.com/minskylab/supersense"

// Store represents an interface where you can preserve your global state
type Store interface {
	CurrentSharedState(lasts int64) (*SharedState, error)
	AddEventToSharedState(event *supersense.Event) error

	SaveCredential(username, password string) error
	UsernameExists(username string) (bool, error)
	ValidateCredential(username, password string) (bool, error)
	UpdateCredential(username, password, newPassword string) error
	ForceUpdateCredential(username, newPassword string) error

	Close() error
}
