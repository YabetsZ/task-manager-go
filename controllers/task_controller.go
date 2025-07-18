package controllers

import (
	"net/http"
	"task-manager/data"
	"task-manager/models"

	"github.com/gin-gonic/gin"
)

// TaskController handles the HTTP requests related to tasks.
type TaskController struct {
	taskService *data.TaskService
}

// NewTaskController is the constructor for TaskController.
func NewTaskController(service *data.TaskService) *TaskController {
	return &TaskController{
		taskService: service,
	}
}

// CreateTask handles POST /tasks requests.
func (tc *TaskController) CreateTask(c *gin.Context) {
	var newTask models.Task

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Basic validation
	if newTask.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	createdTask, err := tc.taskService.CreateTask(newTask)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, createdTask)
}

// GetTasks handles GET /tasks requests.
func (tc *TaskController) GetTasks(c *gin.Context) {
	tasks, err := tc.taskService.GetTasks()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.IndentedJSON(http.StatusOK, tasks)
}

// GetTaskByID handles GET /tasks/:id requests.
func (tc *TaskController) GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	task, err := tc.taskService.GetTaskByID(id)
	if err != nil {
		// Task not found
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, task)
}

// UpdateTask handles PUT /tasks/:id requests.
func (tc *TaskController) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	task, err := tc.taskService.UpdateTask(id, updatedTask)
	if err != nil {
		// Task not found
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

// DeleteTask handles DELETE /tasks/:id requests.
func (tc *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := tc.taskService.DeleteTask(id)
	if err != nil {
		// Task not found
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
