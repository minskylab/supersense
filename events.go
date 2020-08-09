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
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Entities   Entities  `json:"entities"`
	ShareURL   string    `json:"shareUrl"`
	CreatedAt  time.Time `json:"createdAt"`
	EmmitedAt  time.Time `json:"emmitedAt"`
	Message    string    `json:"message"`
	Actor      Person    `json:"actor"`
	SourceID   string    `json:"sourceId"`
	SourceName string    `json:"sourceName"`
	EventKind  string    `json:"eventKind"`
}
