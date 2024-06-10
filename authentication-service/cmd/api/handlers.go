package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	log.Println("Request Payload", requestPayload)

	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	log.Println("User", user)

	if err != nil {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	if user.Active == 0 {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)

	if err != nil || !valid {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	err = app.logRequest("authenticate", requestPayload.Email)

	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Logged in user: " + user.Email,
		Data:    user,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, err := json.MarshalIndent(entry, "", "\t")

	if err != nil {
		return err
	}

	response, err := http.Post("http://logger-service/log", "application/json", bytes.NewBuffer(jsonData))

	if err != nil || response.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
