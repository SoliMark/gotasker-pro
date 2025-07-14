package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/SoliMark/gotasker-pro/config"
	"github.com/SoliMark/gotasker-pro/internal/db"
	"github.com/SoliMark/gotasker-pro/internal/handler"
	"github.com/SoliMark/gotasker-pro/internal/middleware"
	"github.com/SoliMark/gotasker-pro/internal/repository"
	"github.com/SoliMark/gotasker-pro/internal/router"
	"github.com/SoliMark/gotasker-pro/internal/service"
	"github.com/SoliMark/gotasker-pro/internal/util"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Init DB
	dbConn, err := db.NewDB(cfg.DBURL)
	if err != nil {
		log.Fatalf("cannot connect to DB: %v", err)
	}

	// Init JWT Maker
	jwtMaker := util.NewJWTMaker(cfg.JWTSecret)
	if err != nil {
		log.Fatalf("cannot create JWT maker: %v", err)
	}

	userRepo := repository.NewUserRepository(dbConn)
	userService := service.NewUserService(userRepo, jwtMaker)
	userHandler := handler.NewUserHandler(userService)
	authMiddleware := middleware.JWTAuthMiddleware(jwtMaker)

	r := gin.Default()

	router.SetupRoutes(r, userHandler, authMiddleware)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
