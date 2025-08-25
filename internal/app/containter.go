package app

import (
	redis "github.com/redis/go-redis/v9"
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
	RedisClient   *redis.Client
	JWTMiddleware middleware.JWTMiddleware
	UserHandler   *handler.UserHandler
	TaskHandler   *handler.TaskHandler
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

	// Init Redis (optional)
	var redisClient *redis.Client
	if cfg.RedisEnabled() {
		redisClient = db.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	}

	// Init Repository → Service → Handler
	userRepo := repository.NewUserRepository(dbConn)
	userService := service.NewUserService(userRepo, jwtMaker)
	userHandler := handler.NewUserHandler(userService)

	// Init Task components
	taskRepo := repository.NewTaskRepository(dbConn)
	taskService := service.NewTaskService(taskRepo, redisClient, cfg.CacheTTLTasks)
	taskHandler := handler.NewTaskHandler(taskService)

	return &Container{
		Config:        cfg,
		DB:            dbConn,
		RedisClient:   redisClient,
		JWTMiddleware: jwtMiddleware,
		UserHandler:   userHandler,
		TaskHandler:   taskHandler,
	}, nil
}
