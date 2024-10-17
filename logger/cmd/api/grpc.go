package main

import (
	"context"
	"logger/data"
	"logger/logs"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	_, err := logEntry.Insert()

	var res *logs.LogResponse
	if err != nil {
		res = &logs.LogResponse{Result: "failed to logg"}

		return res, err
	}

	res = &logs.LogResponse{Result: "logged!"}

	return res, nil
}
