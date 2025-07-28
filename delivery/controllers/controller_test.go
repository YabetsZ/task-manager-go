package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task-manager/delivery/controllers"
	"task-manager/domain"
	"task-manager/errs"
	"task-manager/usecases/mocks"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ControllerTestSuite struct {
	suite.Suite
	mockTaskUsecase *mocks.TaskUsecase
	mockUserUsecase *mocks.UserUsecase
	controller      *controllers.AppController
	router          *gin.Engine
}

func (s *ControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.mockTaskUsecase = new(mocks.TaskUsecase)
	s.mockUserUsecase = new(mocks.UserUsecase)
	s.controller = controllers.NewAppController(s.mockTaskUsecase, s.mockUserUsecase)

	s.router = gin.Default()
}

func TestAppController(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}

func (s *ControllerTestSuite) performRequest(method, path string, body []byte) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)
	return w
}

// task handler tests

func (s *ControllerTestSuite) TestCreateTask_Success() {
	s.router.POST("/tasks", s.controller.CreateTask)

	dueDate := time.Now()
	taskToCreate := &domain.Task{Title: "Test Task", Description: "something", DueDate: dueDate, Status: "Pending"}
	createdTask := &domain.Task{ID: "123", Title: "Test Task", Description: "something", DueDate: dueDate, Status: "Pending"}

	requestBody, _ := json.Marshal(taskToCreate)

	s.mockTaskUsecase.On("CreateTask", mock.AnythingOfType("*domain.Task")).Return(createdTask, nil).Once()

	w := s.performRequest(http.MethodPost, "/tasks", requestBody)

	s.Require().Equal(http.StatusCreated, w.Code)

	var responseTask domain.Task
	err := json.Unmarshal(w.Body.Bytes(), &responseTask)
	s.Require().NoError(err)
	s.Assert().Equal(createdTask.ID, responseTask.ID)
	s.Assert().Equal(createdTask.Title, responseTask.Title)
	s.Assert().Equal(createdTask.Description, responseTask.Description)
	s.Assert().Equal(createdTask.DueDate.UnixNano(), responseTask.DueDate.UnixNano())
	s.Assert().Equal(createdTask.Status, responseTask.Status)
	s.mockTaskUsecase.AssertExpectations(s.T())
}

func (s *ControllerTestSuite) TestCreateTask_BindingError() {
	s.router.POST("/tasks", s.controller.CreateTask)
	invalidRequestBody := []byte(`{"title": 123}`)

	w := s.performRequest(http.MethodPost, "/tasks", invalidRequestBody)

	s.Assert().Equal(http.StatusBadRequest, w.Code)

	s.mockTaskUsecase.AssertNotCalled(s.T(), "CreateTask")
}

func (s *ControllerTestSuite) TestGetTasks_Success() {
	s.router.GET("/tasks", s.controller.GetTasks)
	mockTasks := []*domain.Task{
		{ID: "1", Title: "Task One"},
		{ID: "2", Title: "Task Two"},
	}
	s.mockTaskUsecase.On("GetTasks").Return(mockTasks, nil).Once()

	w := s.performRequest(http.MethodGet, "/tasks", nil)

	s.Require().Equal(http.StatusOK, w.Code)

	var responseTasks []*domain.Task
	err := json.Unmarshal(w.Body.Bytes(), &responseTasks)
	s.Require().NoError(err)
	s.Assert().Len(responseTasks, 2)
	s.Assert().Equal(mockTasks[0].ID, responseTasks[0].ID)
	s.mockTaskUsecase.AssertExpectations(s.T())
}

func (s *ControllerTestSuite) TestGetTaskByID_NotFound() {
	taskID := "nonexistent"
	s.router.GET("/tasks/:id", s.controller.GetTaskByID)

	s.mockTaskUsecase.On("GetTaskByID", taskID).Return(nil, errs.ErrTaskNotFound).Once()

	w := s.performRequest(http.MethodGet, "/tasks/"+taskID, nil)

	s.Require().Equal(http.StatusNotFound, w.Code)

	var errorResponse map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	s.Require().NoError(err)
	s.Assert().Equal(errs.ErrTaskNotFound.Error(), errorResponse["error"])
	s.mockTaskUsecase.AssertExpectations(s.T())
}

func (s *ControllerTestSuite) TestUpdateTask_Success() {
	s.router.PUT("/tasks/:id", s.controller.UpdateTask)
	taskID := "task123"
	updatePayload := domain.Task{Title: "Updated Title"}
	requestBody, _ := json.Marshal(updatePayload)

	s.mockTaskUsecase.On("UpdateTask", taskID, updatePayload).Return(&updatePayload, nil).Once()

	w := s.performRequest(http.MethodPut, "/tasks/"+taskID, requestBody)

	s.Require().Equal(http.StatusOK, w.Code)
	s.mockTaskUsecase.AssertExpectations(s.T())
}

func (s *ControllerTestSuite) TestUpdateTask_NotFound() {
	s.router.PUT("/tasks/:id", s.controller.UpdateTask)
	taskID := "nonexistent"
	updatePayload := domain.Task{Title: "Updated Title"}
	requestBody, _ := json.Marshal(updatePayload)

	s.mockTaskUsecase.On("UpdateTask", taskID, updatePayload).Return(nil, errs.ErrTaskNotFound).Once()

	w := s.performRequest(http.MethodPut, "/tasks/"+taskID, requestBody)

	s.Assert().Equal(http.StatusNotFound, w.Code)
	s.mockTaskUsecase.AssertExpectations(s.T())
}

func (s *ControllerTestSuite) TestDeleteTask_Success() {
	taskID := "taskToDelete"
	s.router.DELETE("/tasks/:id", s.controller.DeleteTask)
	s.mockTaskUsecase.On("DeleteTask", taskID).Return(nil).Once()

	w := s.performRequest(http.MethodDelete, "/tasks/"+taskID, nil)

	s.Assert().Equal(http.StatusNoContent, w.Code)
	s.mockTaskUsecase.AssertExpectations(s.T())
}

// user handler tests

func (s *ControllerTestSuite) TestLogin_Success() {
	s.router.POST("/login", s.controller.Login)
	loginCreds := map[string]string{"username": "testuser", "password": "password"}
	requestBody, _ := json.Marshal(loginCreds)
	expectedToken := "a.valid.jwt"

	s.mockUserUsecase.On("Login", "testuser", "password").Return(expectedToken, nil).Once()

	w := s.performRequest(http.MethodPost, "/login", requestBody)

	s.Require().Equal(http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.Assert().Equal(expectedToken, response["token"])
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *ControllerTestSuite) TestLogin_InvalidCredentials() {
	s.router.POST("/login", s.controller.Login)
	loginCreds := map[string]string{"username": "testuser", "password": "wrongpassword"}
	requestBody, _ := json.Marshal(loginCreds)

	s.mockUserUsecase.On("Login", "testuser", "wrongpassword").Return("", errs.ErrIncorrectPassword).Once()

	w := s.performRequest(http.MethodPost, "/login", requestBody)

	s.Assert().Equal(http.StatusUnauthorized, w.Code)
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *ControllerTestSuite) TestRegister_Success() {
	s.router.POST("/register", s.controller.Register)
	userPayload := gin.H{"username": "newuser", "password": "password123"}
	requestBody, _ := json.Marshal(userPayload)

	s.mockUserUsecase.On("Register", mock.AnythingOfType("*domain.User")).Return(nil).Once()

	w := s.performRequest(http.MethodPost, "/register", requestBody)

	s.Require().Equal(http.StatusCreated, w.Code)
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *ControllerTestSuite) TestRegister_UsernameExists() {
	s.router.POST("/register", s.controller.Register)
	userPayload := gin.H{"username": "existinguser", "password": "password123"}
	requestBody, _ := json.Marshal(userPayload)

	s.mockUserUsecase.On("Register", mock.AnythingOfType("*domain.User")).Return(errs.ErrUsernameExists).Once()

	w := s.performRequest(http.MethodPost, "/register", requestBody)

	s.Assert().Equal(http.StatusConflict, w.Code)
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *ControllerTestSuite) TestPromote_Success() {
	s.router.POST("/promote/:id", s.controller.Promote)
	userID := "user123"

	s.mockUserUsecase.On("Promote", userID).Return(nil).Once()

	w := s.performRequest(http.MethodPost, "/promote/"+userID, nil)

	s.Require().Equal(http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.Assert().Equal("User promoted to admin successfully", response["message"])
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *ControllerTestSuite) TestPromote_UserNotFound() {
	s.router.POST("/promote/:id", s.controller.Promote)
	userID := "nonexistent"

	s.mockUserUsecase.On("Promote", userID).Return(errs.ErrUserNotFound).Once()

	w := s.performRequest(http.MethodPost, "/promote/"+userID, nil)

	s.Assert().Equal(http.StatusNotFound, w.Code)
	s.mockUserUsecase.AssertExpectations(s.T())
}
