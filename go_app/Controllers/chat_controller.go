package controllers

import (
	"chat_with_go/models"
	"database/sql"
	"log"
)

func ProcessChat(db *sql.DB, request models.ChatCreationRequest) {
	go CreateChat(db, request)
}

func CreateChat(db *sql.DB, request models.ChatCreationRequest) {
	_, err := db.Exec("INSERT INTO chats (application_token, chat_number, name) VALUES (?, ?, ?)", request.ApplicationToken, request.Chat_number, request.Name)
	if err != nil {
		log.Printf("Error creating chat for application token %s: %v", request.ApplicationToken, err)
	} else {
		log.Printf("Chat created for application token: %s", request.ApplicationToken)
	}
}
