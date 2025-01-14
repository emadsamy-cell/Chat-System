package consumer

import (
	controllers "chat_with_go/Controllers"
	"chat_with_go/config"
	"chat_with_go/models"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartChatConsumer(db *sql.DB, ch *amqp.Channel, workers int) {
	msgs, err := ch.Consume(config.ChatQueue, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	messageChannel := make(chan models.ChatCreationRequest, 1000)
	batchChannel := make(chan []models.ChatCreationRequest, 10)

	maxBatchSize := 500
	maxWaitTime := 5 * time.Second

	go func() {
		var messages []models.ChatCreationRequest
		timer := time.NewTimer(maxWaitTime)

		for {
			select {
			case msg := <-messageChannel:
				//log.Printf("Received chat: %v", msg)
				messages = append(messages, msg)
				if len(messages) >= maxBatchSize {
					batchChannel <- messages
					messages = nil
					timer.Reset(maxWaitTime)
				}
			case <-timer.C:
				//log.Printf("Timer expired, sending batch of %d chats", len(messages))
				if len(messages) > 0 {
					batchChannel <- messages
					messages = nil
				}
				timer.Reset(maxWaitTime)
			}
		}
	}()
	// Goroutine to process batches
	go func() {
		for batch := range batchChannel {
			//log.Printf("Processing batch of %d chats", len(batch))
			go controllers.ProcessChat(db, batch)
		}
	}()

	for i := 0; i < workers; i++ {
		go func() {
			for d := range msgs {
				var request models.ChatCreationRequest
				if err := json.Unmarshal(d.Body, &request); err != nil {
					log.Printf("Error decoding chat creation request: %s", err)
					continue
				}
				messageChannel <- request
			}
		}()
	}
}
