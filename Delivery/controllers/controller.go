package controllers

import (
	"net/http"
	"task-manager/data"
	"task-manager/models"

	"github.com/gin-gonic/gin"
)

// AppController handles the HTTP requests in the app.
type AppController struct {
	taskService *data.TaskService
	userService *data.UserService
}

// NewAppController is the constructor for AppController.
func NewAppController(taskService *data.TaskService, userService *data.UserService) *AppController {
	return &AppController{
		taskService, userService,
	}
}

// Task Handler

// CreateTask handles POST api/tasks requests.
func (ac *AppController) CreateTask(c *gin.Context) {
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

	createdTask, err := ac.taskService.CreateTask(newTask)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, createdTask)
}

// GetTasks handles GET api/tasks requests.
func (ac *AppController) GetTasks(c *gin.Context) {
	tasks, err := ac.taskService.GetTasks()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)
}

// GetTaskByID handles GET api/tasks/:id requests.
func (ac *AppController) GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	task, err := ac.taskService.GetTaskByID(id)
	if err != nil {
		// Task not found
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, task)
}

// UpdateTask handles PUT api/tasks/:id requests.
func (ac *AppController) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	task, err := ac.taskService.UpdateTask(id, updatedTask)
	if err != nil {
		// Task not found
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

// DeleteTask handles DELETE api/tasks/:id requests.
func (ac *AppController) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := ac.taskService.DeleteTask(id)
	if err != nil {
		// Task not found
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// User Handlers

// Register handles POST /register requests.
func (ac *AppController) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	err := ac.userService.RegisterUser(&user)
	if err != nil {
		c.JSON(err.Code, gin.H{"error": err.Msg})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login handles POST /login requests.
func (ac *AppController) Login(c *gin.Context) {
	var creds models.User
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	token, err := ac.userService.Login(creds.Username, creds.Password)
	if err != nil {
		c.JSON(err.Code, gin.H{"error": err.Msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Promote handles POST /api/promote requests.
func (ac *AppController) Promote(c *gin.Context) {
	userID := c.Param("id")
	err := ac.userService.Promote(userID)
	if err != nil {
		c.JSON(err.Code, gin.H{"error": err.Msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin successfully"})
}
