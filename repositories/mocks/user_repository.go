package mocks

import (
	"task-manager/domain"

	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepository) GetByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepository) GetByID(id string) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepository) UpdateUserStatus(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *UserRepository) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *UserRepository) CheckUsername(username string) (exist bool, err error) {
	args := m.Called(username)
	return args.Get(0).(bool), args.Error(1)
}
