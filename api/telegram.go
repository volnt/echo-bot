package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joeshaw/envdecode"
	"log"
	"net/http"
	"os"
)

func Webhooks(w http.ResponseWriter, r *http.Request) {
	// Create the Telegram client
	var conf TelegramConfig
	err := envdecode.Decode(&conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "missing env var: [%v]", err)
		return
	}
	tc := NewTelegramClient(conf)

	// Decode the webhook
	var tu TelegramUpdate
	err = json.NewDecoder(r.Body).Decode(&tu)
	if err != nil {
		log.Println("could not decode webhook", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("received webhook", tu)

	// Echo back the message to the user
	err = tc.SendMessage(tu.Message.From.ID, tu.Message.Text)
	if err != nil {
		log.Println("could not send message", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type TelegramConfig struct {
	Token string `env:"TELEGRAM_TOKEN,required"`
}

type TelegramClient interface {
	SendMessage(chatID int, text string) error
}

type TelegramClientImpl struct {
	config TelegramConfig
}

func NewTelegramClient(config TelegramConfig) TelegramClient {
	return &TelegramClientImpl{config}
}

func (c *TelegramClientImpl) SendMessage(chatID int, text string) error {
	tm := TelegramMessage{
		ChatID: chatID,
		Text:   text,
	}

	b, err := json.Marshal(tm)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+c.config.Token+"/sendMessage", bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(req)
	return err
}

type TelegramMessage struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

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
