package usecases_test

import (
	"errors"
	"task-manager/domain"
	"task-manager/errs"
	"task-manager/repositories/mocks"
	"task-manager/usecases"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TaskUsecaseTestSuite struct {
	suite.Suite
	mockTaskRepo *mocks.TaskRepository
	taskUsecase  usecases.TaskUsecase
}

func (s *TaskUsecaseTestSuite) SetupTest() {
	s.mockTaskRepo = new(mocks.TaskRepository)
	s.taskUsecase = usecases.NewTaskUsecase(s.mockTaskRepo)
}

func TestTaskUsecase(t *testing.T) {
	suite.Run(t, new(TaskUsecaseTestSuite))
}

func (s *TaskUsecaseTestSuite) TestCreateTask_Success() {

	inputTask := &domain.Task{Title: "New Task"}
	s.mockTaskRepo.On("Create", inputTask).Return(inputTask, nil).Once()

	createdTask, err := s.taskUsecase.CreateTask(inputTask)

	s.Require().NoError(err)
	s.Assert().Equal(inputTask, createdTask)
	s.mockTaskRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseTestSuite) TestGetTasks_Success() {

	expectedTasks := []*domain.Task{
		{ID: "1", Title: "Task 1"},
		{ID: "2", Title: "Task 2"},
	}
	s.mockTaskRepo.On("GetAll").Return(expectedTasks, nil).Once()

	tasks, err := s.taskUsecase.GetTasks()

	s.Require().NoError(err)
	s.Assert().Len(tasks, 2)
	s.Assert().Equal(expectedTasks, tasks)
	s.mockTaskRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseTestSuite) TestGetTaskByID_Success() {

	taskID := "task123"
	expectedTask := &domain.Task{ID: taskID, Title: "Found Task"}
	s.mockTaskRepo.On("GetByID", taskID).Return(expectedTask, nil).Once()

	task, err := s.taskUsecase.GetTaskByID(taskID)

	s.Require().NoError(err)
	s.Assert().Equal(expectedTask, task)
	s.mockTaskRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseTestSuite) TestGetTaskByID_NotFound() {

	taskID := "nonexistent"
	s.mockTaskRepo.On("GetByID", taskID).Return(nil, errs.ErrTaskNotFound).Once()

	task, err := s.taskUsecase.GetTaskByID(taskID)

	s.Require().Error(err)
	s.Assert().ErrorIs(err, errs.ErrTaskNotFound, "Expected a specific task not found error")
	s.Assert().Nil(task)
	s.mockTaskRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseTestSuite) TestUpdateTask_Success() {

	taskID := "task123"
	taskUpdate := domain.Task{Title: "Updated Title"}
	expectedUpdatedTask := &domain.Task{ID: taskID, Title: "Updated Title", Status: domain.StatusPending}

	s.mockTaskRepo.On("Update", taskID, taskUpdate).Return(expectedUpdatedTask, nil).Once()

	updatedTask, err := s.taskUsecase.UpdateTask(taskID, taskUpdate)

	s.Require().NoError(err)
	s.Assert().Equal(expectedUpdatedTask, updatedTask)
	s.mockTaskRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseTestSuite) TestDeleteTask_Success() {

	taskID := "task123"
	s.mockTaskRepo.On("Delete", taskID).Return(nil).Once()

	err := s.taskUsecase.DeleteTask(taskID)

	s.Require().NoError(err)
	s.mockTaskRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseTestSuite) TestDeleteTask_Failure() {

	taskID := "task123"
	expectedErr := errors.New("database error")
	s.mockTaskRepo.On("Delete", taskID).Return(expectedErr).Once()

	err := s.taskUsecase.DeleteTask(taskID)

	s.Require().Error(err)
	s.Assert().Equal(expectedErr, err)
	s.mockTaskRepo.AssertExpectations(s.T())
}
