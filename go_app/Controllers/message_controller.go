package controllers

import (
	"chat_with_go/models"
	"chat_with_go/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func ProcessMessage(db *sql.DB, request models.MessageCreationRequest) {
	go CreateMessage(db, request)

	go IndexMessage(request)
}

func CreateMessage(db *sql.DB, request models.MessageCreationRequest) {
	_, err := db.Exec("INSERT INTO messages (application_token, chat_number, body, message_number) VALUES (?, ?, ?, ?)",
		request.ApplicationToken, request.Chat_number, request.Body, request.Message_number)
	if err != nil {
		log.Printf("Error creating message for chat number %d: %v", request.Chat_number, err)
	} else {
		log.Printf("Message created for chat number %d", request.Chat_number)
	}
}

func IndexMessage(request models.MessageCreationRequest) {
	es := utils.GetElasticsearchClient()
	docID := fmt.Sprintf("%s_%d_%d", request.ApplicationToken, request.Chat_number, request.Message_number)

	message := map[string]interface{}{
		"application_token": request.ApplicationToken,
		"chat_number":       request.Chat_number,
		"body":              request.Body,
		"message_number":    request.Message_number,
	}

	body, _ := json.Marshal(message)

	res, err := es.Index(
		"messages",
		strings.NewReader(string(body)),
		es.Index.WithDocumentID(docID),
	)

	if err != nil {
		log.Printf("Error indexing message for chat Number %d: %v", request.Chat_number, err)
	} else {
		log.Printf("Message indexed for chat Number %d", request.Chat_number)
	}

	defer res.Body.Close()
}
