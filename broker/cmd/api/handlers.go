package main

import (
	"broker/event"
	"broker/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"net/rpc"
	"time"
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
		app.logItemViaRPC(w, requestPayload.Log)
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

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)

	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Logged !",
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEmitter(app.Rabbit)

	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	jsonData, _ := json.Marshal(payload)

	err = emitter.Push(string(jsonData), "log.INFO")
	return err

}

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) logItemViaRPC(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")

	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}
	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string

	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: result,
	}

	_ = app.writeJSON(w, http.StatusAccepted, response)
}

func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	_ = app.readJSON(w, r, &requestPayload)

	conn, err := grpc.NewClient("logger-service:50001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{MinConnectTimeout: 5 * time.Second}),
	)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	client := logs.NewLogServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.WriteLog(ctx, &logs.LogRequest{LogEntry: &logs.Log{
		Name: requestPayload.Log.Name,
		Data: requestPayload.Log.Data,
	}})

	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	jsonresponse := jsonResponse{
		Error:   false,
		Message: response.Result,
	}

	_ = app.writeJSON(w, http.StatusAccepted, jsonresponse)
}
