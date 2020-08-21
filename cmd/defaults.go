package main

import (
	"time"

	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/config"
	"github.com/minskylab/supersense/sources"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func defaultSources(conf *config.Config) ([]supersense.Source, error) {
	defaultSources := make([]supersense.Source, 0)

	if conf.DummyPeriod != "" && conf.DummyMessage != "" {
		dur, err := time.ParseDuration(conf.DummyPeriod)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		dummy, err := sources.NewDummy(dur, conf.DummyMessage)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		defaultSources = append(defaultSources, dummy)
		log.WithFields(log.Fields{
			"period":  dur,
			"message": conf.DummyMessage,
		}).Info("Dummy source activated")
	}

	if len(conf.GithubRepos) != 0 {
		var token *string = nil
		if conf.GithubToken != "" {
			token = &conf.GithubToken
		}

		github, err := sources.NewGithub(token, conf.GithubRepos)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		defaultSources = append(defaultSources, github)
		log.WithFields(log.Fields{
			"repos":     conf.GithubRepos,
			"withToken": token != nil,
		}).Info("Github source activated")
	}

	if conf.TwitterAccessSecret != "" && conf.TwitterAccessToken != "" &&
		conf.TwitterConsumerKey != "" && conf.TwitterConsumerSecret != "" &&
		len(conf.TwitterQuery) != 0 {

		twitter, err := sources.NewTwitter(sources.TwitterClientProps{
			ConsumerKey:    conf.TwitterConsumerKey,
			ConsumerSecret: conf.TwitterConsumerSecret,
			AccessToken:    conf.TwitterAccessToken,
			AccessSecret:   conf.TwitterAccessSecret,
			QueryToTrack:   conf.TwitterQuery,
		})
		if err != nil {
			return nil, errors.WithStack(err)
		}

		defaultSources = append(defaultSources, twitter)
		log.WithFields(log.Fields{
			"query": conf.TwitterQuery,
		}).Info("Twitter source activated")
	}

	return defaultSources, nil
}
