package main

import (
	"context"
	"strings"

	"github.com/joho/godotenv"
	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/config"
	"github.com/minskylab/supersense/server"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func launchDefaultService() error {
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

	eventsPipe := make(chan *supersense.Event, 1)
	done := make(chan struct{})

	go mux.Register(eventsPipe, done)

	go func(eventsPipe chan *supersense.Event) {
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
				"[%s] %s: %s | by: %s @%s",
				event.SourceName,
				event.Title,
				cutMessage,
				event.Actor.Name,
				username,
			)
		}
	}(eventsPipe)

	ctx := context.Background()
	if err := mux.RunAllSources(ctx); err != nil {
		return errors.WithStack(err)
	}

	return server.LaunchServer(mux, conf.Port, conf.GraphQLPlayground)
}
