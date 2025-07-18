package main

import (
	"context"
	"log"
	"task-manager/controllers"
	"task-manager/data"
	"task-manager/router"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGO_URI       = "mongodb://root:password@localhost:27017"
	DATABASE_NAME   = "task_db"
	COLLECTION_NAME = "tasks"
)

func main() {
	client, err := connectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	collection := client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)
	newTaskService := data.NewTaskService(collection)

	newTaskController := controllers.NewTaskController(newTaskService)
	r := router.SetupRouter(newTaskController)

	log.Println("Starting server on port 5000...")
	if err := r.Run(":5000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func connectToDB() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		return nil, err
	}

	// Ping the primary to verify the connection.
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MongoDB!")
	return client, nil
}
