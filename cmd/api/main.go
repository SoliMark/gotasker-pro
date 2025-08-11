package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/SoliMark/gotasker-pro/internal/app"
	"github.com/SoliMark/gotasker-pro/internal/router"
)

func main() {
	container, err := app.InitApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	r := gin.Default()
	router.SetupRoutes(r, container)

	if err := r.Run(":" + container.Config.AppPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
