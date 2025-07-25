package infrastructure

import (
	"fmt"
	"task-manager/errs"
	"task-manager/usecases"

	"golang.org/x/crypto/bcrypt"
)

type bcryptService struct{}

func NewBcryptService() usecases.PasswordService {
	return &bcryptService{}
}

func (s *bcryptService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%w: %v", errs.ErrUnexpected, err)
	}
	return string(bytes), nil
}

func (s *bcryptService) Compare(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errs.ErrIncorrectPassword
	}
	return nil
}
