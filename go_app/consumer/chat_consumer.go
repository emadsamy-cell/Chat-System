package consumer

import (
	controllers "chat_with_go/Controllers"
	"chat_with_go/config"
	"chat_with_go/models"
	"database/sql"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartChatConsumer(db *sql.DB, ch *amqp.Channel, workers int) {
	msgs, err := ch.Consume(config.ChatQueue, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < workers; i++ {
		go func() {
			for d := range msgs {
				var request models.ChatCreationRequest
				if err := json.Unmarshal(d.Body, &request); err != nil {
					log.Printf("Error decoding chat creation request: %s", err)
					continue
				}
				controllers.ProcessChat(db, request)
			}
		}()
	}
}
