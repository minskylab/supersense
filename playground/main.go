package main

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	"github.com/minskylab/supersense/sources"
)

func main() {
	godotenv.Load() // loading .env vars
	authToken := os.Getenv("GITHUB_TOKEN")

	source, err := sources.NewGithub([]string{
		"minskylab/supersense",
	}, &authToken)
	if err != nil {
		log.Panic(err)
	}

	source.Run(context.TODO())

	// channel := source.Events(context.TODO())

	time.Sleep(1 * time.Hour)
}
