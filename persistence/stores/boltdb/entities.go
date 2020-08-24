package boltdb

import (
	"time"

	"github.com/minskylab/supersense"
)

// Event describes a simple event from a source
type Event struct {
	supersense.Event `storm:"inline"`
	ID string `storm:"id"`
	EmittedAt time.Time `storm:"index"`
}
