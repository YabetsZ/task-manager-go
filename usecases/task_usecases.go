package usecases

import (
	"task-manager/domain"
)

type TaskUsecase interface {
	CreateTask(task *domain.Task) (*domain.Task, error)
	GetTasks() ([]*domain.Task, error)
	GetTaskByID(id string) (*domain.Task, error)
	UpdateTask(id string, updatedTask domain.Task) (*domain.Task, error)
	DeleteTask(id string) error
}

// TaskRepository defines the interface for task data operations.
type TaskRepository interface {
	Create(task *domain.Task) (*domain.Task, error)
	GetAll() ([]*domain.Task, error)
	GetByID(id string) (*domain.Task, error)
	Update(id string, updatedTask domain.Task) (*domain.Task, error)
	Delete(id string) error
}

type taskUsecase struct {
	taskRepo TaskRepository
}

func NewTaskUsecase(ur TaskRepository) TaskUsecase {
	return &taskUsecase{
		taskRepo: ur,
	}
}

func (ts *taskUsecase) CreateTask(task *domain.Task) (*domain.Task, error) {
	return ts.taskRepo.Create(task)
}

func (ts *taskUsecase) GetTasks() ([]*domain.Task, error) {

	return ts.taskRepo.GetAll()
}

func (ts *taskUsecase) GetTaskByID(id string) (*domain.Task, error) {
	return ts.taskRepo.GetByID(id)
}

func (ts *taskUsecase) UpdateTask(id string, updatedTask domain.Task) (*domain.Task, error) {

	return ts.taskRepo.Update(id, updatedTask)
}

func (ts *taskUsecase) DeleteTask(id string) error {
	return ts.taskRepo.Delete(id)
}
