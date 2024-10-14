package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}
type MailPayload struct {
	From    string `json:"from,omitempty"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.log(w, requestPayload.Log)
	case "mail":
		app.mail(w, requestPayload.Mail)

	default:
		_ = app.errorJSON(w, errors.New("unknown action"))
	}

}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.Marshal(a)

	req, err := http.NewRequest("POST", "http://authentication-service:8080/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		_ = app.errorJSON(w, errors.New("invalid credentials"))
		return
	}

	var jsonFromAuth jsonResponse

	err = json.NewDecoder(res.Body).Decode(&jsonFromAuth)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	if jsonFromAuth.Error {
		_ = app.errorJSON(w, errors.New(fmt.Sprint(jsonFromAuth.Message)), http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated !",
		Data:    jsonFromAuth.Data,
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) log(w http.ResponseWriter, l LogPayload) {
	jsonData, _ := json.Marshal(l)

	request, err := http.NewRequest("POST", "http://logger-service:8080/log", bytes.NewBuffer(jsonData))
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		_ = app.errorJSON(w, errors.New("log did not created "))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Logged !",
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) mail(w http.ResponseWriter, m MailPayload) {
	jsonData, _ := json.Marshal(m)

	req, err := http.NewRequest("POST", "http://mail-service:8080/send", bytes.NewBuffer(jsonData))

	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		_ = app.errorJSON(w, errors.New("mail did not sent"))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "mail sent",
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}
