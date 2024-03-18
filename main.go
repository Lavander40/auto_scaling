package main

import (
	tgclient "auto_scaling/clients/telegram"
	"auto_scaling/consumer/event_consumer.go"
	"auto_scaling/events/telegram"
	"auto_scaling/scaler/yandex"
	"auto_scaling/storage/sqlite"
	"context"
	"flag"
	"log"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	fileStoragePath   = "data/user_storage"
	batchSize         = 100
)

func main() {
	// st := files.New(fileStoragePath)
	st, err := sqlite.New(context.TODO(), sqliteStoragePath)
	if err != nil {
		log.Fatalf("can't connect to storage: ", err)
	}

	if err := st.Init(); err != nil {
		log.Fatalf("can't init storage: ", err)
	}

	processor := telegram.New(
		tgclient.New(tgBotHost, mustToken()),
		st,
		yandex.New(),
	)

	log.Print("service started")

	// consumer(fetcher, processor)
	consumer := event_consumer.New(processor, processor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String("token", "", "telegramm api token")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is unset")
	}

	return *token
}
