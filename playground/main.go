package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	"github.com/minskylab/supersense/sources"
)

func main() {
	godotenv.Load() // loading .env vars
	source, err := sources.NewGithub([]string{
		"minskylab/supersense",
	})
	if err != nil {
		log.Panic(err)
	}

	source.Run(context.TODO())

	time.Sleep(1 * time.Hour)
}
