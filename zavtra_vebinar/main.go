package main

import (
	"log"

	tgClient "ZavtraVebinar/clients/telegram"
	event_consumer "ZavtraVebinar/consumer/event-consumer"
	"ZavtraVebinar/events/telegram"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)

func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
	)

	log.Print("Сервис запущен")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("Сервис остановлен", err)
	}
}

func mustToken() string {
	token := "5006959166:AAFrCXi3qnRphI3HdwMWgTBxYJrN3kQHh1g"

	return token
}
