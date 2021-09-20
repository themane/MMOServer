package mongoRepository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	connectTimeoutSecs = 100
)

func getConnection(mongoURL string, ctx context.Context) *mongo.Client {
	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			log.Print(evt.Command)
		},
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL).SetMonitor(cmdMonitor))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to MongoDB")
	return client
}

func disconnect(client *mongo.Client, ctx context.Context) {
	err := client.Disconnect(ctx)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Println("Connection to MongoDB closed.")
}
