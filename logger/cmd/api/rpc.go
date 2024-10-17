package main

import (
	"log"
	"logger/data"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	logEntry := data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	}

	_, err := logEntry.Insert()
	if err != nil {
		log.Println(err)
		return err
	}

	*resp = "Processed payload via RPC" + payload.Name
	return nil

}
