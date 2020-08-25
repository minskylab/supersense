package persistence

import (
	"time"

	"github.com/minskylab/supersense"
)

// SharedConfig is useful to change the shared config params for all your spectators
type SharedConfig struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// SharedState refers to an event board shared for all of your observers (consumers)
type SharedState struct {
	Board      []*supersense.Event `json:"board"`
	LastUpdate time.Time           `json:"last_update"`
}
