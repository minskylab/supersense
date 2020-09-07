package sources

import (
	"time"

	"github.com/minskylab/supersense"
	uuid "github.com/satori/go.uuid"
)

type Spokesman struct {
	id       string
	name     string
	username string
	email    string
	channel  chan supersense.Event
	sourceName string
}

func NewSpokesman(name, username, email string) (*Spokesman, error) {
	source := &Spokesman{
		id:       uuid.NewV4().String(),
		sourceName: "spokesman",
		name:     name,
		username: username,
		email:    email,
		channel:  make(chan supersense.Event),
	}
	return source, nil
}

// Identify implements the Source interface
func (s *Spokesman) Identify(nameOrID string) bool {
	return s.sourceName == nameOrID || s.id == nameOrID
}

func (s *Spokesman) Run() error {
	return nil
}

func (s *Spokesman) Dispose() {
	close(s.channel)
}

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

func (s *Spokesman) Broadcast(title, message string, entities supersense.Entities, url string, delay *time.Duration) {
	go s.programBroadcast(nil, nil, nil, title, message, entities, url, delay)
}

func (s *Spokesman) BroadcastWithActor(name string, username, photo *string, title, message string, entities supersense.Entities, url string, delay *time.Duration) {
	go s.programBroadcast(&name, username, photo, title, message, entities, url, delay)
}