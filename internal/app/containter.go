package app

import (
	"gorm.io/gorm"

	"github.com/SoliMark/gotasker-pro/config"
	"github.com/SoliMark/gotasker-pro/internal/db"
	"github.com/SoliMark/gotasker-pro/internal/handler"
	"github.com/SoliMark/gotasker-pro/internal/middleware"
	"github.com/SoliMark/gotasker-pro/internal/repository"
	"github.com/SoliMark/gotasker-pro/internal/service"
	"github.com/SoliMark/gotasker-pro/internal/util"
)

type Container struct {
	Config        *config.Config
	DB            *gorm.DB
	JWTMiddleware middleware.JWTMiddleware
	UserHandler   *handler.UserHandler

	// TODO:
	// RedisClient   *redis.Client
	// TaskHandler   *handler.TaskHandler
}

func InitApp() (*Container, error) {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Init DB
	dbConn, err := db.NewDB(cfg.DBURL)
	if err != nil {
		return nil, err
	}

	// Init JWT
	jwtMaker := util.NewJWTMaker(cfg.JWTSecret)
	jwtMiddleware := middleware.JWTAuthMiddleware(jwtMaker)

	// Init Repository → Service → Handler
	userRepo := repository.NewUserRepository(dbConn)
	userService := service.NewUserService(userRepo, jwtMaker)
	userHandler := handler.NewUserHandler(userService)

	return &Container{
		Config:        cfg,
		DB:            dbConn,
		JWTMiddleware: jwtMiddleware,
		UserHandler:   userHandler,
	}, nil
}
