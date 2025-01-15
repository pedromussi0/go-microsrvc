package main

import (
	"context"
	"log"
	"log-service/data"
	"time"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	// Call the logger service
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.Background(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo", err)
	}

	*resp = "Processed Payload via RPC" + payload.Name

	return nil
}
