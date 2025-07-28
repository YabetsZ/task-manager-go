package mocks

import (
	"task-manager/domain"

	"github.com/stretchr/testify/mock"
)

type TaskUsecase struct {
	mock.Mock
}

func (m *TaskUsecase) CreateTask(task *domain.Task) (*domain.Task, error) {
	args := m.Called(task)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}
func (m *TaskUsecase) GetTasks() ([]*domain.Task, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Task), args.Error(1)
}
func (m *TaskUsecase) GetTaskByID(id string) (*domain.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}
func (m *TaskUsecase) UpdateTask(id string, task domain.Task) (*domain.Task, error) {
	args := m.Called(id, task)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}
func (m *TaskUsecase) DeleteTask(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
