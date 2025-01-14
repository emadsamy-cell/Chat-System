package controllers

import (
	"bytes"
	"chat_with_go/models"
	"chat_with_go/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

func ProcessMessage(db *sql.DB, requests []models.MessageCreationRequest) {
	go CreateMessage(db, requests)

	go IndexMessage(requests)
}

func CreateMessage(db *sql.DB, requests []models.MessageCreationRequest) {
	if len(requests) == 0 {
		return
	}

	query := "INSERT INTO messages (application_token, chat_number, body, message_number) VALUES "
	vals := []interface{}{}

	for i, req := range requests {
		if i > 0 {
			query += ","
		}

		query += "(?, ?, ?, ?)"
		vals = append(vals, req.ApplicationToken, req.Chat_number, req.Body, req.Message_number)
	}

	_, err := db.Exec(query, vals...)
	if err != nil {
		log.Printf("Error executing bulk message creation: %v", err)
		return
	}

	log.Printf("Successfully created %d messages in batch", len(requests))
}

func IndexMessage(requests []models.MessageCreationRequest) {
	es := utils.GetElasticsearchClient()
	var buf bytes.Buffer

	for _, request := range requests {
		docID := fmt.Sprintf("%s_%d_%d", request.ApplicationToken, request.Chat_number, request.Message_number)
		meta := []byte(fmt.Sprintf(`{ "index" : { "_index" : "messages", "_id" : "%s" } }%s`, docID, "\n"))

		message := map[string]interface{}{
			"application_token": request.ApplicationToken,
			"chat_number":       request.Chat_number,
			"body":              request.Body,
			"message_number":    request.Message_number,
		}

		data, err := json.Marshal(message)
		if err != nil {
			log.Fatalf("Cannot encode message %d: %s", request.Message_number, err)
		}
		data = append(data, "\n"...)

		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)
	}

	res, err := es.Bulk(bytes.NewReader(buf.Bytes()), es.Bulk.WithIndex("messages"))

	if err != nil {
		log.Printf("Bulk indexing error: %s", err)
	} else {
		log.Printf("Bulk indexing successful, length: %d", len(requests))
	}

	res.Body.Close()
	buf.Reset()
}
