package boltdb

import "github.com/minskylab/supersense"

// Event describes a simple event from a source
type Event struct {
	supersense.Event `storm:"inline"`
	ID string `storm:"id"`
	EmittedAt string `storm:"index"`
}
