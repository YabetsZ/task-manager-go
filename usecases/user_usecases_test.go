package usecases_test

import (
	"task-manager/domain"
	"task-manager/errs"
	"task-manager/infrastructure"
	"task-manager/repositories/mocks"
	"task-manager/usecases"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserUsecaseTestSuite struct {
	suite.Suite
	mockUserRepo *mocks.UserRepository
	// TODO: In a full test suite, these would also be mocks.
	passwordService usecases.PasswordService
	jwtService      usecases.JWTService
	userUsecase     usecases.UserUsecase
}

func (s *UserUsecaseTestSuite) SetupTest() {
	s.mockUserRepo = new(mocks.UserRepository)
	s.passwordService = infrastructure.NewBcryptService()
	s.jwtService = infrastructure.NewJWTServiceV5()
	s.userUsecase = usecases.NewUserUsecase(s.mockUserRepo, s.passwordService, s.jwtService)
}

func TestUserUsecase(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}

// Test cases

func (s *UserUsecaseTestSuite) TestRegister_Success_FirstUserIsAdmin() {
	// Arrange
	user := &domain.User{
		Username:     "admin",
		PasswordHash: "password123",
	}
	// Use mock.Anything because we don't care about the context value in this test.
	s.mockUserRepo.On("CheckUsername", user.Username).Return(false, nil)
	s.mockUserRepo.On("Count").Return(int64(0), nil).Once()
	s.mockUserRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil).Once()

	// Act
	err := s.userUsecase.Register(user)

	// Assert
	s.Require().NoError(err)
	s.Assert().Equal(domain.RoleAdmin, user.Role, "The first user should be an admin")

	// Verify that the mock methods were called as expected
	s.mockUserRepo.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestRegister_Success_SecondUserIsRegularUser() {
	user := &domain.User{Username: "testuser", PasswordHash: "password123"}
	// Arrange
	s.mockUserRepo.On("Count").Return(int64(1), nil).Once()
	s.mockUserRepo.On("CheckUsername", user.Username).Return(false, nil)
	s.mockUserRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil).Once()

	// Act
	err := s.userUsecase.Register(user)

	// Assert
	s.Require().NoError(err)
	s.Assert().Equal(domain.RoleUser, user.Role, "Subsequent users should have the 'user' role")
	s.mockUserRepo.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestLogin_Failure_UserNotFound() {
	// We configure the mock to return our specific application error
	s.mockUserRepo.On("GetByUsername", "nonexistent").Return(nil, errs.ErrUserNotFound).Once()

	// Act
	token, err := s.userUsecase.Login("nonexistent", "password")

	// Assert
	s.Require().Error(err, "Expected an error for non-existent user")
	s.Assert().Equal(errs.ErrUserNotFound, err, "Error should be ErrInvalidUserId")
	s.Assert().Empty(token, "Token should be empty on failure")
	s.mockUserRepo.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestLogin_Failure_WrongPassword() {
	// Hashed version of "correctpassword"
	hashedPassword, _ := s.passwordService.Hash("correctpassword")
	mockUser := &domain.User{Username: "testuser", PasswordHash: hashedPassword}

	s.mockUserRepo.On("GetByUsername", "testuser").Return(mockUser, nil).Once()

	// Act
	token, err := s.userUsecase.Login("testuser", "wrongpassword")

	// Assert
	s.Require().Error(err, "Expected an error for wrong password")
	s.Assert().Equal(errs.ErrIncorrectPassword, err, "Error should be ErrIncorrectPassword")
	s.Assert().Empty(token, "Token should be empty on failure")
	s.mockUserRepo.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestLogin_Success() {
	hashedPassword, _ := s.passwordService.Hash("correctpassword")
	mockUser := &domain.User{Username: "testuser", PasswordHash: hashedPassword}

	s.mockUserRepo.On("GetByUsername", "testuser").Return(mockUser, nil).Once()

	// Act
	token, err := s.userUsecase.Login("testuser", "correctpassword")

	// Assert
	s.Require().NoError(err)
	s.Assert().NotEmpty(token, "A JWT token should be returned on successful login")
	s.mockUserRepo.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestGetUserByID_Success() {
	// Arrange
	expectedUserID := "some_valid_id"
	mockUser := &domain.User{Username: "testuser", Role: domain.RoleUser}

	s.mockUserRepo.On("GetByID", expectedUserID).Return(mockUser, nil).Once()

	// Act
	user, err := s.userUsecase.GetUserByID(expectedUserID)

	// Assert
	s.Require().NoError(err, "GetUserByID should not return an error on success")
	s.Assert().NotNil(user, "Returned user should not be nil")
	s.Assert().Equal(mockUser.Username, user.Username, "Returned user should match the mock")
	s.mockUserRepo.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestGetUserByID_NotFound() {
	// Arrange
	nonExistentID := "non_existent_id"

	// We expect GetByID to be called, and we instruct the mock to return our application's not-found error.
	s.mockUserRepo.On("GetByID", nonExistentID).Return(nil, errs.ErrInvalidUserId).Once()

	// Act
	user, err := s.userUsecase.GetUserByID(nonExistentID)

	// Assert
	s.Require().Error(err, "GetUserByID should return an error when user is not found")
	s.Assert().ErrorIs(err, errs.ErrInvalidUserId, "The error should be a specific ErrInvalidUserId")
	s.Assert().Nil(user, "The returned user should be nil on failure")
	s.mockUserRepo.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestPromote_Success() {
	// Arrange
	userIDToPromote := "user_to_promote_id"

	s.mockUserRepo.On("UpdateUserStatus", userIDToPromote).Return(nil).Once()

	// Act
	err := s.userUsecase.Promote(userIDToPromote)

	// Assert
	s.Require().NoError(err, "Promote should not return an error on success")
	s.mockUserRepo.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestPromote_UserNotFound() {
	// Arrange
	nonExistentID := "non_existent_id"
	s.mockUserRepo.On("UpdateUserStatus", nonExistentID).Return(errs.ErrUserNotFound).Once()

	// Act
	err := s.userUsecase.Promote(nonExistentID)

	// Assert
	s.Require().Error(err, "Promote should return an error if the user is not found")
	s.Assert().ErrorIs(err, errs.ErrUserNotFound, "The error should be ErrUserNotFound")

	s.mockUserRepo.AssertExpectations(s.T())
}
