package router

import (
	"github.com/SoliMark/gotasker-pro/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userHandler *handler.UserHandler, authMiddleware gin.HandlerFunc) {
	// Public routes
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// Protected routes with JWT
	api := r.Group("/api")
	api.Use(authMiddleware)
	{
		// User Profile
		api.GET("/profile", userHandler.Profile)
		// Future: Task CRUD example
		//api.GET("/task",taskHandler.ListTask)
	}

}
