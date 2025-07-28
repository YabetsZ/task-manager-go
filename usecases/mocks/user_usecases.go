package mocks

import (
	"task-manager/domain"

	"github.com/stretchr/testify/mock"
)

type UserUsecase struct {
	mock.Mock
}

func (m *UserUsecase) Register(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *UserUsecase) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}
func (m *UserUsecase) Promote(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *UserUsecase) GetUserByID(id string) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
