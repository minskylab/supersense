package main

import (
	"context"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/sources"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	godotenv.Load() // loading .env vars

	ctx := context.TODO()

	dummySource, err := sources.NewDummy(5*time.Minute, "Hello World")
	if err != nil {
		log.Panic(err)
	}

	authToken := os.Getenv("GITHUB_TOKEN")

	githubSource, err := sources.NewGithub([]string{"minskylab/supersense"}, &authToken)
	if err != nil {
		log.Panic(err)
	}

	twitterSource, err := sources.NewTwitter(sources.TwitterClientProps{
		ConsumerKey:    os.Getenv("CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("CONSUMER_SECRET"),
		AccessToken:    os.Getenv("ACCESS_TOKEN"),
		AccessSecret:   os.Getenv("ACCESS_SECRET"),
		QueryToTrack:   []string{"#peru"},
	})
	if err != nil {
		log.Panic(err)
	}

	mux, err := supersense.NewMux(ctx, twitterSource, dummySource, githubSource)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		for event := range mux.Events() {
			log.Infof(spew.Sdump(event))
		}
	}()

	mux.RunAllSources(ctx)

	time.Sleep(10 * time.Hour)
}
