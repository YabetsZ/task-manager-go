package router

import (
	"task-manager/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(tc *controllers.TaskController) *gin.Engine {
	r := gin.Default()

	taskRoutes := r.Group("/tasks")

	taskRoutes.POST("", tc.CreateTask)
	taskRoutes.GET("", tc.GetTasks)
	taskRoutes.GET("/:id", tc.GetTaskByID)
	taskRoutes.PUT("/:id", tc.UpdateTask)
	taskRoutes.DELETE("/:id", tc.DeleteTask)

	return r
}
