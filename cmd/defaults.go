package main

import (
	"strings"
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

		// Preprocessing and cleaning TwitterQuery
		// Why?. Answer: https://github.com/moby/moby/issues/20169

		repos := make([]string, 0)
		for _, r := range conf.GithubRepos {
			repo := strings.ReplaceAll(r, "\"", "")
			repo = strings.TrimSpace(repo)
			repos = append(repos, repo)
		}

		github, err := sources.NewGithub(token, conf.GithubRepos)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		defaultSources = append(defaultSources, github)
		log.WithFields(log.Fields{
			"repos":     repos,
			"withToken": token != nil,
		}).Info("Github source activated")
	}

	if conf.TwitterAccessSecret != "" && conf.TwitterAccessToken != "" &&
		conf.TwitterConsumerKey != "" && conf.TwitterConsumerSecret != "" &&
		len(conf.TwitterQuery) != 0 {

		// Preprocessing and cleaning TwitterQuery
		// Why?. Answer: https://github.com/moby/moby/issues/20169

		query := make([]string, 0)
		for _, q := range conf.TwitterQuery {
			q1 := strings.ReplaceAll(q, "\"", "")
			q1 = strings.TrimSpace(q)
			query = append(query, q1)
		}

		twitter, err := sources.NewTwitter(sources.TwitterClientProps{
			ConsumerKey:    conf.TwitterConsumerKey,
			ConsumerSecret: conf.TwitterConsumerSecret,
			AccessToken:    conf.TwitterAccessToken,
			AccessSecret:   conf.TwitterAccessSecret,
			QueryToTrack:   query,
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

func specialSpokesman(conf *config.Config) (*sources.Spokesman, error){
	var spokesman *sources.Spokesman
	var err error


	if conf.Spokesman == true && conf.SpokesmanName != "" && conf.SpokesmanUsername != "" {
		spokesman, err = sources.NewSpokesman(conf.SpokesmanName, conf.SpokesmanUsername, conf.SpokesmanEmail)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		log.WithFields(log.Fields{
			"name":     conf.SpokesmanName,
			"username": conf.SpokesmanUsername,
			"email":    conf.SpokesmanEmail,
		}).Info("Spokesman source activated")
	}

	return spokesman,nil
}