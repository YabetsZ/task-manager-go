package data

import (
	"fmt"
	"task-manager/models"

	"github.com/google/uuid"
)

type TaskService struct {
	tasks map[string]*models.Task
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks: make(map[string]*models.Task),
	}
}

func (ts *TaskService) CreateTask(task models.Task) *models.Task {
	if task.Status == "" {
		task.Status = models.StatusPending
	}

	task.ID = uuid.NewString()

	ts.tasks[task.ID] = &task

	return &task
}

func (ts *TaskService) GetTasks() []*models.Task {
	tasks := make([]*models.Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

func (ts *TaskService) GetTaskByID(id string) (*models.Task, error) {
	task, exist := ts.tasks[id]
	if !exist {
		return nil, fmt.Errorf("couldn't find a taks with id %s", id)
	}
	return task, nil
}

func (ts *TaskService) UpdateTask(id string, updatedTask models.Task) (*models.Task, error) {
	if _, exist := ts.tasks[id]; !exist {
		return nil, fmt.Errorf("couldn't find a taks with id %s", id)
	}
	if updatedTask.Description != "" {
		ts.tasks[id].Description = updatedTask.Description
	}
	if updatedTask.DueDate.IsZero() {
		ts.tasks[id].DueDate = updatedTask.DueDate
	}
	if updatedTask.Status != "" {
		ts.tasks[id].Status = updatedTask.Status
	}
	if updatedTask.Title != "" {
		ts.tasks[id].Title = updatedTask.Title
	}

	return ts.tasks[id], nil
}

func (ts *TaskService) DeleteTask(id string) error {
	if _, err := ts.GetTaskByID(id); err != nil {
		return err
	}
	delete(ts.tasks, id)
	return nil
}
