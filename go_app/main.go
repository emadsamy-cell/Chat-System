package main

import (
	"chat_with_go/consumer"
	"chat_with_go/jobs"
	"chat_with_go/utils"
	"log"
	"runtime"
	"time"
)

func main() {
	db, err := utils.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(runtime.NumCPU() * 4)
	db.SetMaxIdleConns(runtime.NumCPU())
	db.SetConnMaxLifetime(5 * time.Minute)

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
