package main

import (
	"net/http"
)

type mailMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) sendMail(w http.ResponseWriter, r *http.Request) {
	var requestPayload mailMessage

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSTMPMessage(&msg)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Message sent to " + requestPayload.To,
	}

	if err := app.writeJSON(w, http.StatusOK, payload); err != nil {
		app.errorJSON(w, err)
	}
}
