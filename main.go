package main

import (
	"log"
	"task-manager/controllers"
	"task-manager/data"
	"task-manager/router"
)

func main() {
	newTaskService := data.NewTaskService()
	newTaskController := controllers.NewTaskController(newTaskService)
	r := router.SetupRouter(newTaskController)

	log.Println("Starting server on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
