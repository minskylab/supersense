package main

import (
	"time"

	"github.com/minskylab/supersense/sources"
	log "github.com/sirupsen/logrus"
)

func main() {
	source, err := sources.NewDummy(2*time.Second, "Hello World")
	if err != nil {
		log.Panic(err)
	}

	go func() {
		for event := range *source.Events() {
			log.Infof("%s | %s", event.ID, event.Message)
		}
	}()

	source.Run()

	time.Sleep(10 * time.Second)
}
