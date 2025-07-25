package repositories

import (
	"context"
	"fmt"
	"log"
	"task-manager/domain"
	"task-manager/errs"
	"task-manager/usecases"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoTaskRepository struct {
	collection *mongo.Collection
}

func NewMongoTaskRepository(collection *mongo.Collection) usecases.TaskRepository {
	return &mongoTaskRepository{
		collection,
	}
}

type mongoTask struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	DueDate     time.Time          `bson:"due_date"`
	Status      string             `bson:"status"`
}

func (t *mongoTaskRepository) buildTask(from mongoTask) (to *domain.Task) {
	return &domain.Task{
		ID:          from.ID.Hex(),
		Title:       from.Title,
		Description: from.Description,
		DueDate:     from.DueDate,
		Status:      from.Status,
	}
}

func (t *mongoTaskRepository) Create(task *domain.Task) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mTask := mongoTask{
		ID:          primitive.NewObjectID(),
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Status:      "",
	}

	if task.Status != domain.StatusCompleted && task.Status != domain.StatusInProgress {
		mTask.Status = domain.StatusPending
	} else {
		mTask.Status = task.Status
	}

	_, err := t.collection.InsertOne(ctx, mTask)
	if err != nil {
		return nil, err
	}

	return t.buildTask(mTask), nil
}

func (t *mongoTaskRepository) GetAll() ([]*domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("[mongoTaskRepository] GetAll: fetching all tasks")

	tasks := make([]*domain.Task, 0)

	cursor, err := t.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("[mongoTaskRepository] GetAll: Find error %e", err)
		return nil, fmt.Errorf("%w: %v", errs.ErrUnexpected, err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task mongoTask
		if err := cursor.Decode(&task); err != nil {
			log.Printf("Decode error: %v", err)
			return nil, fmt.Errorf("%w: %v", errs.ErrUnexpected, err)
		}
		tasks = append(tasks, t.buildTask(task))
	}

	log.Printf("[mongoTaskRepository] GetAll: retrieved tasks %d", len(tasks))
	return tasks, nil
}

func (t *mongoTaskRepository) GetByID(id string) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.ErrInvalidTaskId
	}

	var task mongoTask

	err = t.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.ErrTaskNotFound
		}
		return nil, fmt.Errorf("%w: %v", errs.ErrUnexpected, err)
	}
	return t.buildTask(task), nil
}
func (t *mongoTaskRepository) Update(id string, updatedTask domain.Task) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.ErrInvalidTaskId
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

	res, err := t.collection.UpdateByID(ctx, objID, update)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrUnexpected, err)
	}
	if res.MatchedCount == 0 {
		return nil, errs.ErrInvalidTaskId
	}

	return t.GetByID(id)
}

func (t *mongoTaskRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errs.ErrInvalidTaskId
	}

	res, err := t.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errs.ErrTaskNotFound
	}
	return nil
}
