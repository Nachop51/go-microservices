package main

import (
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

	log.Println("there is a user", user)

	if user.Active == 0 {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	log.Println("is user active", user.Active)

	log.Println("User", user)
	log.Println("User", requestPayload.Password)

	valid, err := user.PasswordMatches(requestPayload.Password)

	if err != nil || !valid {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Logged in user: " + user.Email,
		Data:    user,
	}

	app.writeJSON(w, http.StatusOK, payload)
}
