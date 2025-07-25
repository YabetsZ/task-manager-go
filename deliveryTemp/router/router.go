package router

import (
	"task-manager/delivery/controllers"
	"task-manager/domain"
	"task-manager/infrastructure"
	"task-manager/usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter(ac *controllers.AppController, uu usecases.UserUsecase) *gin.Engine {
	r := gin.Default()

	// public routes
	r.POST("/register", ac.Register)
	r.POST("/login", ac.Login)

	// private routes
	api := r.Group("/api")
	{
		// Admin-only routes
		adminRoutes := api.Group("")
		adminRoutes.Use(infrastructure.AuthMiddleware(uu, domain.RoleAdmin))
		{
			adminRoutes.POST("/tasks", ac.CreateTask)
			adminRoutes.PUT("/tasks/:id", ac.UpdateTask)
			adminRoutes.DELETE("/tasks/:id", ac.DeleteTask)
			adminRoutes.POST("/promote/:id", ac.Promote)
		}

		// Routes for all authenticated users (Admin and User)
		userRoutes := api.Group("")
		userRoutes.Use(infrastructure.AuthMiddleware(uu, domain.RoleUser))
		{
			userRoutes.GET("/tasks", ac.GetTasks)
			userRoutes.GET("/tasks/:id", ac.GetTaskByID)
		}
	}

	return r
}
