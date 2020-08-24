package supersense

import "time"

// MediaEntity represents a media in the event (e.g. Photos, Videos, Files)
type MediaEntity struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

// URLEntity wrap a simple URL from the event message
type URLEntity struct {
	URL        string `json:"url"`
	DisplayURL string `json:"display_url"`
}

// Entities saves three types of entities
type Entities struct {
	Tags  []string      `json:"tags"`
	Media []MediaEntity `json:"media"`
	Urls  []URLEntity   `json:"urls"`
}

// Event describes a simple event from a source
type Event struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Entities   Entities  `json:"entities"`
	ShareURL   string    `json:"shareURL"`
	CreatedAt  time.Time `json:"createdAt"`
	EmittedAt  time.Time `json:"emittedAt"`
	Message    string    `json:"message"`
	Actor      Person    `json:"actor"`
	SourceID   string    `json:"sourceId"`
	SourceName string    `json:"sourceName"`
	EventKind  string    `json:"eventKind"`
	// New
	Labels []string `json:"labels"`
}
