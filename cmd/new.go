package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/config"
	"github.com/minskylab/supersense/persistence"
	"github.com/minskylab/supersense/persistence/stores/boltdb"
	"github.com/minskylab/supersense/server"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func internalSubscriber(eventsPipe chan *supersense.Event, store persistence.Store) {
	for event := range eventsPipe {
		maxLength := 20
		cutMessage := event.Message
		if len(event.Message) > maxLength {
			cutMessage = cutMessage[:maxLength]
			cutMessage = strings.Replace(cutMessage, "\n", " ", -1)
			cutMessage = strings.Trim(cutMessage, "\n ")
		}

		username := ""
		if event.Actor.Username != nil {
			username = *event.Actor.Username
		}

		log.Debugf(
			"[%s] %s | by: %s @%s",
			event.SourceName,
			cutMessage,
			event.Actor.Name,
			username,
		)

		if store != nil {
			if err := store.AddEventToSharedState(event); err != nil {
				log.Error(fmt.Sprintf("%+v", err))
			}
		}
	}
}

func setupPersistence(conf *config.Config) (persistence.Store, error) {
	var store persistence.Store
	var err error

	if conf.Persistence {
		if conf.PersistenceRedisAddress != "" && conf.PersistenceRedisPassword != "" {
			// TODO: Implement redis store
			log.Warn("At the moment the redis persistence layer are not implemented :(")
		} else if conf.PersistenceBoltDBFilePath != "" {
			store, err = boltdb.NewStore(conf.PersistenceBoltDBFilePath, -1)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			log.WithField("filepath", conf.PersistenceBoltDBFilePath).Info("Persistence layer activated with boltdb")
		}

	}

	return store, nil
}

func setupCredentials(conf config.Config, store persistence.Store) error {
	username := conf.RootCredentialUsername
	if len(username) < 4 {
		return errors.New("invalid username, please choose other")
	}

	password := conf.RootCredentialPassword

	if password == "" {
		log.Info("Root password not detected")

		charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		source := rand.New(rand.NewSource(time.Now().Unix()))

		b := make([]byte, 12)

		for i := range b {
			b[i] = charset[source.Intn(len(charset))]
		}

		password = string(b)
		log.Info("New auto-generated password: " + password)
	}

	valid, err := store.ValidateCredential(username, password)
	if err != nil {
		return errors.WithStack(err)
	}

	if !valid {
		if err := store.SaveCredential(username, password); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func launchDefaultService(done chan struct{}) error {
	_ = godotenv.Load() // loading .env vars

	conf, err := config.LoadDefault()
	if err != nil {
		return errors.WithStack(err)
	}

	if conf.Debug {
		log.SetLevel(log.DebugLevel)
	}

	loadedSources, err := defaultSources(conf)
	if err != nil {
		return errors.WithStack(err)
	}

	mux, err := supersense.NewMux(loadedSources...)
	if err != nil {
		return errors.WithStack(err)
	}

	spokesman, err := specialSpokesman(conf)
	if err != nil {
		return errors.WithStack(err)
	}

	if spokesman != nil {
		mux.AddNewSource(spokesman)
	}

	store, err := setupPersistence(conf)
	if err != nil {
		return errors.WithStack(err)
	}

	eventsPipe := make(chan *supersense.Event, 1)
	go mux.Register(eventsPipe, done)
	go internalSubscriber(eventsPipe, store)

	if err := mux.RunAllSources(); err != nil {
		return errors.WithStack(err)
	}

	// TODO add functional options as a struct builder

	return server.LaunchServer(mux, conf.Port, conf.GraphQLPlayground, spokesman, store)
}
