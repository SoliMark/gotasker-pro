package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	rediscontainer "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"

	"github.com/SoliMark/gotasker-pro/internal/app"
	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/router"
	"github.com/SoliMark/gotasker-pro/internal/util"
)

// ContainerTestSuite 整合測試套件
type ContainerTestSuite struct {
	suite.Suite
	postgresContainer *postgrescontainer.PostgresContainer
	redisContainer    *rediscontainer.RedisContainer
	db                *gorm.DB
	app               *app.Container
	router            *gin.Engine
	jwtMaker          *util.JWTMaker
}

// SetupSuite 設置測試套件
func (ts *ContainerTestSuite) SetupSuite() {
	ctx := context.Background()

	// 啟動 PostgreSQL 容器
	postgresContainer, err := postgrescontainer.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgrescontainer.WithDatabase("gotasker_test"),
		postgrescontainer.WithUsername("testuser"),
		postgrescontainer.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(120*time.Second),
		),
	)
	require.NoError(ts.T(), err)
	ts.postgresContainer = postgresContainer

	// 啟動 Redis 容器
	redisContainer, err := rediscontainer.RunContainer(ctx,
		testcontainers.WithImage("redis:7-alpine"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithOccurrence(1).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(ts.T(), err)
	ts.redisContainer = redisContainer

	// 獲取容器連接信息
	postgresHost, err := postgresContainer.Host(ctx)
	require.NoError(ts.T(), err)
	postgresPort, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(ts.T(), err)
	redisHost, err := redisContainer.Host(ctx)
	require.NoError(ts.T(), err)
	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	require.NoError(ts.T(), err)

	// 設置環境變數
	ts.setupEnvironment(postgresHost, postgresPort.Port(), redisHost, redisPort.Port())

	// 初始化應用
	ts.setupApplication()
}

// SetupTest 每個測試前的設置
func (ts *ContainerTestSuite) SetupTest() {
	ts.cleanupDatabase()
}

// TearDownSuite 清理測試套件
func (ts *ContainerTestSuite) TearDownSuite() {
	ctx := context.Background()
	if ts.postgresContainer != nil {
		ts.postgresContainer.Terminate(ctx)
	}
	if ts.redisContainer != nil {
		ts.redisContainer.Terminate(ctx)
	}
}

// setupEnvironment 設置測試環境變數
func (ts *ContainerTestSuite) setupEnvironment(postgresHost, postgresPort, redisHost, redisPort string) {
	// 設置環境變數供應用使用
	ts.T().Setenv("DB_URL", fmt.Sprintf("postgres://testuser:testpass@%s:%s/gotasker_test?sslmode=disable", postgresHost, postgresPort))
	ts.T().Setenv("REDIS_ADDR", fmt.Sprintf("%s:%s", redisHost, redisPort))
	ts.T().Setenv("JWT_SECRET", "test-secret-key-for-integration-tests")
}

// setupApplication 初始化應用
func (ts *ContainerTestSuite) setupApplication() {
	var err error

	// 初始化應用容器
	ts.app, err = app.InitApp()
	require.NoError(ts.T(), err)

	// 獲取數據庫連接
	ts.db = ts.app.DB

	// 創建 JWT Maker 實例
	ts.jwtMaker = util.NewJWTMaker("test-secret-key-for-integration-tests")

	// 設置 Gin 路由
	ts.router = gin.Default()
	router.SetupRoutes(ts.router, ts.app)

	// 執行數據庫遷移
	ts.runMigrations()
}

// runMigrations 執行數據庫遷移
func (ts *ContainerTestSuite) runMigrations() {
	// 執行 GORM 自動遷移
	err := ts.db.AutoMigrate(&model.User{}, &model.Task{})
	if err != nil {
		log.Printf("Migration failed: %v", err)
	}
	log.Println("Database migrations completed")
}

// cleanupDatabase 清理數據庫
func (ts *ContainerTestSuite) cleanupDatabase() {
	ts.db.Exec("DELETE FROM tasks WHERE 1=1")
	ts.db.Exec("DELETE FROM users WHERE 1=1")
	ts.db.Exec("ALTER SEQUENCE IF EXISTS users_id_seq RESTART WITH 1")
	ts.db.Exec("ALTER SEQUENCE IF EXISTS tasks_id_seq RESTART WITH 1")
}

// createTestUserAndLogin 創建測試用戶並登入，返回 token
func (ts *ContainerTestSuite) createTestUserAndLogin(email, password, username string) string {
	// 註冊用戶
	registerBody := map[string]interface{}{
		"username": username,
		"email":    email,
		"password": password,
	}
	registerJSON, _ := json.Marshal(registerBody)
	registerReq, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(registerJSON))
	registerReq.Header.Set("Content-Type", "application/json")
	registerW := httptest.NewRecorder()
	ts.router.ServeHTTP(registerW, registerReq)

	// 檢查註冊是否成功
	if registerW.Code != http.StatusOK {
		return "" // 註冊失敗，返回空字符串
	}

	// 登入用戶
	loginBody := map[string]interface{}{
		"email":    email,
		"password": password,
	}
	loginJSON, _ := json.Marshal(loginBody)
	loginReq, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJSON))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	ts.router.ServeHTTP(loginW, loginReq)

	// 檢查登入是否成功
	if loginW.Code != http.StatusOK {
		return "" // 登入失敗，返回空字符串
	}

	// 解析登入響應獲取 token
	var loginResponse map[string]interface{}
	json.Unmarshal(loginW.Body.Bytes(), &loginResponse)
	if token, exists := loginResponse["token"]; exists && token != nil {
		return token.(string)
	}
	return "" // 返回空字符串如果登入失敗
}
