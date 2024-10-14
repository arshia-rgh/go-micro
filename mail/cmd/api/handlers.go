package main

import (
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from,omitempty"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		_ = app.errorJSON(w, err)
		log.Println(err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		_ = app.errorJSON(w, err)
		log.Println(err)

		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Mail sent to " + requestPayload.To,
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)

}
