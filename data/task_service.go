package data

import (
	"context"
	"fmt"
	"task-manager/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskService struct {
	collection *mongo.Collection
}

func NewTaskService(collection *mongo.Collection) *TaskService {
	return &TaskService{
		collection,
	}
}

func (ts *TaskService) CreateTask(task models.Task) (*models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if task.Status == "" {
		task.Status = models.StatusPending
	}

	task.ID = primitive.NewObjectID()

	_, err := ts.collection.InsertOne(ctx, task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (ts *TaskService) GetTasks() ([]*models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tasks := make([]*models.Task, 0)

	cursor, err := ts.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (ts *TaskService) GetTaskByID(id string) (*models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	var task models.Task

	err = ts.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("couldn't find a taks with id %s", id)
		}
		return nil, err
	}

	return &task, nil
}

func (ts *TaskService) UpdateTask(id string, updatedTask models.Task) (*models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	updateFields := bson.M{}
	if updatedTask.Description != "" {
		updateFields["description"] = updatedTask.Description
	}
	if !updatedTask.DueDate.IsZero() {
		updateFields["due_date"] = updatedTask.DueDate
	}
	if updatedTask.Status != "" {
		updateFields["status"] = updatedTask.Status
	}
	if updatedTask.Title != "" {
		updateFields["title"] = updatedTask.Title
	}
	update := bson.M{
		"$set": updateFields,
	}

	res, err := ts.collection.UpdateByID(ctx, objID, update)
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, fmt.Errorf("couldn't find a taks with id %s", id)
	}

	return ts.GetTaskByID(id)
}

func (ts *TaskService) DeleteTask(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	res, err := ts.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("couldn't find a taks with id %s", id)
	}

	return nil
}
