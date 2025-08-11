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
		// Future: Task CRUD example
		//api.GET("/task",taskHandler.ListTask)
	}

}
