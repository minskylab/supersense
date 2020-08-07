package supersense

import "time"

// Event describes a simple event from a source
type Event struct {
	ID        string
	Title     string
	BannerURL string
	ShareURL  string
	EmmitedAt time.Time
	Message   string
	Person    Person
	SourceID  string
}
