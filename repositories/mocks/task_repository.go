package mocks

import (
	"task-manager/domain"

	"github.com/stretchr/testify/mock"
)

// TaskRepository is a mock type for the TaskRepository interface
type TaskRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: task
func (m *TaskRepository) Create(task *domain.Task) (*domain.Task, error) {
	args := m.Called(task)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

// GetAll provides a mock function
func (m *TaskRepository) GetAll() ([]*domain.Task, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Task), args.Error(1)
}

// GetByID provides a mock function with given fields: id
func (m *TaskRepository) GetByID(id string) (*domain.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

// Update provides a mock function with given fields: id, updatedTask
func (m *TaskRepository) Update(id string, updatedTask domain.Task) (*domain.Task, error) {
	args := m.Called(id, updatedTask)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

// Delete provides a mock function with given fields: id
func (m *TaskRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
