package models

type MessageCreationRequest struct {
	ApplicationToken string `json:"application_token"`
	Chat_number      int    `json:"chat_number"`
	Body             string `json:"body"`
	Message_number   int    `json:"message_number"`
}
