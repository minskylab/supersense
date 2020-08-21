package main

import (
	"context"
	"fmt"
	"strings"

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

		log.Infof(
			"[%s] %s | by: %s @%s",
			event.SourceName,
			cutMessage,
			event.Actor.Name,
			username,
		)

		if store != nil {
			if _, err := store.AddEventToSharedState(event); err != nil {
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
			log.WithField("filepath",conf.PersistenceBoltDBFilePath ).Info("Persistence layer activated with boltdb")
		}

	}

	return store, nil
}

func launchDefaultService(done chan struct{}) error {
	_ = godotenv.Load() // loading .env vars

	conf, err := config.LoadDefault()
	if err != nil {
		return errors.WithStack(err)
	}

	loadedSources, err := defaultSources(conf)
	if err != nil {
		return errors.WithStack(err)
	}

	mux, err := supersense.NewMux(loadedSources...)
	if err != nil {
		return errors.WithStack(err)
	}

	store, err := setupPersistence(conf)
	if err != nil {
		return errors.WithStack(err)
	}


	eventsPipe := make(chan *supersense.Event, 1)
	go mux.Register(eventsPipe, done)
	go internalSubscriber(eventsPipe, store)

	ctx := context.Background()
	if err := mux.RunAllSources(ctx); err != nil {
		return errors.WithStack(err)
	}

	return server.LaunchServer(mux, conf.Port, conf.GraphQLPlayground)
}
