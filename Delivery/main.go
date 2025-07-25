package main

import (
	"context"
	"log"
	"task-manager/delivery/controllers"
	"task-manager/delivery/router"
	"task-manager/infrastructure"
	"task-manager/repositories"
	"task-manager/usecases"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGO_URI     = "mongodb://root:password@localhost:27017"
	DATABASE_NAME = "task_db"
)

func main() {
	client, err := connectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	tasksCollection := client.Database(DATABASE_NAME).Collection("tasks")
	usersCollection := client.Database(DATABASE_NAME).Collection("users")
	newMongoTaskRepository := repositories.NewMongoTaskRepository(tasksCollection)
	newMongoUserRepository := repositories.NewMongoUserRepository(usersCollection)
	newTaskUseCase := usecases.NewTaskUsecase(newMongoTaskRepository)
	newUserUsecase := usecases.NewUserUsecase(
		newMongoUserRepository,
		infrastructure.NewBcryptService(),
		infrastructure.NewJWTServiceV5(),
	)

	newAppController := controllers.NewAppController(newTaskUseCase, newUserUsecase)
	r := router.SetupRouter(newAppController, newUserUsecase)

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

	return client, nil
}
