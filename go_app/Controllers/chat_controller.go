package controllers

import (
	"chat_with_go/models"
	"database/sql"
	"log"
)

func ProcessChat(db *sql.DB, requests []models.ChatCreationRequest) {
	CreateChat(db, requests)
}

func CreateChat(db *sql.DB, requests []models.ChatCreationRequest) {
	if len(requests) == 0 {
		return
	}

	query := "INSERT INTO chats (application_token, chat_number, name) VALUES "
	vals := []interface{}{}

	for i, req := range requests {
		if i > 0 {
			query += ","
		}

		query += "(?, ?, ?)"
		vals = append(vals, req.ApplicationToken, req.Chat_number, req.Name)
	}

	_, err := db.Exec(query, vals...)
	if err != nil {
		log.Printf("Error executing bulk chat creation: %v", err)
		return
	}

	log.Printf("Successfully created %d chats in batch", len(requests))
}
