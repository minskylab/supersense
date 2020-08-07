package supersense

import "time"

// MediaEntity represents a media in the event (e.g. Photos, Videos, Files)
type MediaEntity struct {
	URL  string
	Type string
}

// URLEntity wrap a simple URL from the event message
type URLEntity struct {
	URL        string
	DisplayURL string
}

// Entities saves three types of entities
type Entities struct {
	Tags  []string
	Media []MediaEntity
	Urls  []URLEntity
}

// Event describes a simple event from a source
type Event struct {
	ID         string
	Title      string
	Entities   Entities
	ShareURL   string
	CreatedAt  time.Time
	EmmitedAt  time.Time
	Message    string
	Person     Person
	SourceID   string
	SourceName string
	EventKind  string
}
