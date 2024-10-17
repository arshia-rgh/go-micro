package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"logger/data"
	"logger/logs"
	"net"
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

func (app *Config) gRPCListen() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%v", gRpcPort))

	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Printf("gRPC server started on port %v", gRpcPort)

	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
