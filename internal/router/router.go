package router

import (
	"github.com/gin-gonic/gin"

	"github.com/SoliMark/gotasker-pro/internal/app"
)

func SetupRoutes(r *gin.Engine, c *app.Container) {
	// Public routes
	r.POST("/register", c.UserHandler.Register)
	r.POST("/login", c.UserHandler.Login)

	// Protected routes with JWT
	api := r.Group("/api")
	api.Use(c.JWTMiddleware)
	{
		// User Profile
		api.GET("/profile", c.UserHandler.Profile)

		// Task CRUD
		tasks := api.Group("/tasks")
		{
			tasks.POST("", c.TaskHandler.CreateTask)
			tasks.GET("", c.TaskHandler.ListTasks)
			tasks.GET("/:id", c.TaskHandler.GetTask)
			tasks.PUT("/:id", c.TaskHandler.UpdateTask)
			tasks.DELETE("/:id", c.TaskHandler.DeleteTask)
		}
	}
}
