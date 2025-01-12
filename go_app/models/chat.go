package models

type ChatCreationRequest struct {
	ApplicationToken string `json:"application_token"`
	Chat_number      int    `json:"chat_number"`
	Name             string `json:"name"`
}
