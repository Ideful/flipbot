package main

import (
	"flag"
	"log"

	"github.com/Ideful/flipbot/clients/telegram"
	eventconsumer "github.com/Ideful/flipbot/event-consumer"
	telegram2 "github.com/Ideful/flipbot/events/telegram"
	"github.com/Ideful/flipbot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {

	tgClient := telegram.New(tgBotHost, mustToken())

	eventsProcessor := telegram2.New(tgClient, files.New(storagePath))
	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service stopped")
	}
}

func mustToken() string {
	token := flag.String("token", "", "token for bot access")

	flag.Parse()

	if *token == "" {
		log.Fatal("empty token")
	}
	return *token
}
