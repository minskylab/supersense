package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/sources"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	godotenv.Load() // loading .env vars

	dummySource, err := sources.NewDummy(5*time.Second, "Hello World")
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

	mux, err := supersense.NewMux(dummySource, twitterSource)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		for event := range mux.Events() {
			log.Infof("%s | %s", event.ID, event.Message)
		}
	}()

	mux.RunAllSources()

	time.Sleep(1000 * time.Second)
}
