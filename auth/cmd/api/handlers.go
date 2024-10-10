package main

import (
	"errors"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	type payLoad struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var payload payLoad
	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(payload.Email)

	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordsMatches(payload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: "logged in",
		Data:    user,
	}

	err = app.writeJSON(w, http.StatusAccepted, response)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
}
