package main

import (
	"chat_with_go/consumer"
	"chat_with_go/jobs"
	"chat_with_go/utils"
	"log"
	"runtime"
)

func main() {
	db, err := utils.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	workers := runtime.NumCPU() * 2

	conn, ch, err := utils.SetupRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	defer ch.Close()

	go consumer.StartChatConsumer(db, ch, workers)
	go consumer.StartMessageConsumer(db, ch, workers)

	go jobs.BatchUpdateCounts(db)
	select {}
}
