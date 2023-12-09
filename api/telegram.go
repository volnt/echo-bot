package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type TelegramUpdate struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID int `json:"id"`
		} `json:"from"`
		Text string `json:"text"`
		Date int    `json:"date"`
	} `json:"message"`
}

func Webhooks(w http.ResponseWriter, r *http.Request) {
	var tu TelegramUpdate

	err := json.NewDecoder(r.Body).Decode(&tu)
	if err != nil {
		log.Println("could not decode webhook", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("received webhook", tu)
}
