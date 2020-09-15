package sources

import (
	"time"

	"github.com/minskylab/supersense"
	uuid "github.com/satori/go.uuid"
)

// Spokesman is a source useful to emit custom message to your event board
type Spokesman struct {
	id         string
	name       string
	username   string
	email      string
	channel    chan supersense.Event
	sourceName string
}

// NewSpokesman creates and return a new bootstraped Spokesman struct
func NewSpokesman(name, username, email string) (*Spokesman, error) {
	source := &Spokesman{
		id:         uuid.NewV4().String(),
		sourceName: "spokesman",
		name:       name,
		username:   username,
		email:      email,
		channel:    make(chan supersense.Event),
	}
	return source, nil
}

// Identify implements the Source interface
func (s *Spokesman) Identify(nameOrID string) bool {
	return s.sourceName == nameOrID || s.id == nameOrID
}

// Run starts the spokesman source
func (s *Spokesman) Run() error {
	return nil
}

// Dispose close all streams and flows with the source
func (s *Spokesman) Dispose() {
	close(s.channel)
}

// Pipeline is necessary to retreive the spokesman events chanel
func (s *Spokesman) Pipeline() <-chan supersense.Event {
	return s.channel
}

func (s *Spokesman) programBroadcast(name, username, photo *string, title, message string, entities supersense.Entities, url string, delay *time.Duration) {
	if delay != nil {
		time.Sleep(*delay)
	}

	finalName := s.name
	if name != nil {
		finalName = *name
	}

	finalUsername := s.username
	if username != nil {
		finalUsername = *username
	}

	finalPhoto := "https://api.adorable.io/avatars/72/" + finalUsername + ".png"
	if photo != nil {
		finalPhoto = *photo
	}

	event := supersense.Event{
		ID:        uuid.NewV4().String(),
		Title:     title,
		Entities:  entities,
		ShareURL:  url,
		CreatedAt: time.Now(),
		Message:   message,
		Actor: supersense.Person{
			Name:     finalName,
			Photo:    finalPhoto,
			Owner:    "supersense",
			Email:    &s.email,
			Username: &finalUsername,
		},
		SourceID:   s.id,
		SourceName: "spokesman",
		EventKind:  "broadcast",
	}

	event.EmittedAt = time.Now()
	s.channel <- event
}

// Broadcast emit a new message without external actor information
func (s *Spokesman) Broadcast(title, message string, entities supersense.Entities, url string, delay *time.Duration) {
	go s.programBroadcast(nil, nil, nil, title, message, entities, url, delay)
}

// BroadcastWithActor emit a new message with custom actor information
func (s *Spokesman) BroadcastWithActor(name string, username, photo *string, title, message string, entities supersense.Entities, url string, delay *time.Duration) {
	go s.programBroadcast(&name, username, photo, title, message, entities, url, delay)
}
