package usecases

import (
	"log"
	"task-manager/domain"
	"task-manager/errs"
)

type UserUsecase interface {
	Register(user *domain.User) error
	Login(username, password string) (string, error)
	Promote(userID string) error
	GetUserByID(id string) (*domain.User, error)
}

type JWTService interface {
	GenerateJWT(*domain.User) (string, error)
}

type PasswordService interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	Create(user *domain.User) error
	GetByUsername(username string) (*domain.User, error)
	GetByID(id string) (*domain.User, error)
	UpdateUserStatus(id string) error
	Count() (int64, error)
	CheckUsername(username string) (exist bool, err error)
}

type userUsecase struct {
	userRepo    UserRepository
	passwordSvc PasswordService
	jwtSvc      JWTService
}

func NewUserUsecase(ur UserRepository, ps PasswordService, js JWTService) UserUsecase {
	return &userUsecase{
		userRepo:    ur,
		passwordSvc: ps,
		jwtSvc:      js,
	}
}

func (u *userUsecase) Register(user *domain.User) error {
	// check if username already exists
	exist, err := u.userRepo.CheckUsername(user.Username)
	if err != nil {
		return err
	}
	if exist {
		return errs.ErrUsernameExists
	}

	hashedPassword, err := u.passwordSvc.Hash(user.Password)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword

	// For the requirement: "If the database is empty, the first created user will be an admin."
	count, err := u.userRepo.Count()
	if err != nil {
		return err
	}
	if count == 0 {
		user.Role = domain.RoleAdmin
	} else {
		user.Role = domain.RoleUser
	}

	err = u.userRepo.Create(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) Login(username, password string) (string, error) {

	log.Printf("INFO: Login attempt for username: '%s'", username)

	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		return "", err
	}

	err = u.passwordSvc.Compare(user.PasswordHash, password)
	if err != nil {
		log.Printf("WARN: Login failed for username '%s': invalid password", username)
		return "", err
	}

	log.Printf("INFO: User '%s' (ID: %s, Role: %s) successfully authenticated", user.Username, user.ID, user.Role)

	return u.jwtSvc.GenerateJWT(user)
}

// This function will be called by the auth middleware
func (u *userUsecase) GetUserByID(id string) (*domain.User, error) {

	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) Promote(userID string) error {
	log.Printf("Attempting to promote user in usecase, %v", userID)
	err := u.userRepo.UpdateUserStatus(userID)
	if err != nil {
		return err
	}

	return nil
}
