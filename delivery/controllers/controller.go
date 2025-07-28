package controllers

import (
	"log"
	"net/http"
	"task-manager/domain"
	"task-manager/usecases"
	"time"

	"github.com/gin-gonic/gin"
)

// AppController handles the HTTP requests in the app.
type AppController struct {
	taskUsecase usecases.TaskUsecase
	userUsecase usecases.UserUsecase
}

type ginTask struct {
	ID          string    `json:"id,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

func fromDomainTask(task *domain.Task) *ginTask {
	return &ginTask{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Status:      task.Status,
	}
}
func toDomainTask(gtask *ginTask) *domain.Task {
	return &domain.Task{
		ID:          gtask.ID,
		Title:       gtask.Title,
		Description: gtask.Description,
		DueDate:     gtask.DueDate,
		Status:      gtask.Status,
	}
}

func NewAppController(tu usecases.TaskUsecase, uu usecases.UserUsecase) *AppController {
	return &AppController{taskUsecase: tu, userUsecase: uu}
}

// Task Handler

// CreateTask handles POST api/tasks requests.
func (ac *AppController) CreateTask(c *gin.Context) {

	var newTask ginTask

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Basic validation
	if newTask.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	createdTask, err := ac.taskUsecase.CreateTask(toDomainTask(&newTask))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdTask)
}

// GetTasks handles GET api/tasks requests.
func (ac *AppController) GetTasks(c *gin.Context) {
	tasks, err := ac.taskUsecase.GetTasks()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ginTasks := make([]*ginTask, 0)
	for _, task := range tasks {
		ginTasks = append(ginTasks, fromDomainTask(task))
	}

	c.IndentedJSON(http.StatusOK, ginTasks)
}

// GetTaskByID handles GET api/tasks/:id requests.
func (ac *AppController) GetTaskByID(c *gin.Context) {

	id := c.Param("id")
	task, err := ac.taskUsecase.GetTaskByID(id)
	if err != nil {
		// Task not found
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, fromDomainTask(task))
}

// UpdateTask handles PUT api/tasks/:id requests.
func (ac *AppController) UpdateTask(c *gin.Context) {

	id := c.Param("id")
	var updatedTask ginTask
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	task, err := ac.taskUsecase.UpdateTask(id, *toDomainTask(&updatedTask))
	if err != nil {
		// Task not found
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fromDomainTask(task))
}

// DeleteTask handles DELETE api/tasks/:id requests.
func (ac *AppController) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := ac.taskUsecase.DeleteTask(id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// User Handlers

type ginUser struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func fromDomainUser(user *domain.User) *ginUser {
	return &ginUser{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}
}
func toDomainUser(gUser *ginUser) *domain.User {
	return &domain.User{
		ID:       gUser.ID,
		Username: gUser.Username,
		Password: gUser.Password,
		Role:     gUser.Password,
	}
}

// Register handles POST /register requests.
func (ac *AppController) Register(c *gin.Context) {

	var user ginUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	err := ac.userUsecase.Register(toDomainUser(&user))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login handles POST /login requests.
func (ac *AppController) Login(c *gin.Context) {
	var creds ginUser
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	token, err := ac.userUsecase.Login(creds.Username, creds.Password)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Promote handles POST /api/promote requests.
func (ac *AppController) Promote(c *gin.Context) {
	userID := c.Param("id")
	log.Println("Attempting to promote user", userID)
	err := ac.userUsecase.Promote(userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin successfully"})
}
