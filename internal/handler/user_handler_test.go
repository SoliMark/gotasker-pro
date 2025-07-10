package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SoliMark/gotasker-pro/internal/constant"
	"github.com/SoliMark/gotasker-pro/internal/handler"
	"github.com/SoliMark/gotasker-pro/internal/service"
	"github.com/SoliMark/gotasker-pro/internal/service/mock_service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockUserService(ctrl)

	userHandler := &handler.UserHandler{
		UserService: mockSvc,
	}

	body := handler.RegisterRequest{
		Email:    "test@example.com",
		Password: "plaintextpassword",
	}

	jsonBody, _ := json.Marshal(body)

	mockSvc.
		EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(nil)

	router := gin.Default()
	router.POST("/register", userHandler.Register)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "User registered successfully")
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockUserService(ctrl)

	userHandler := &handler.UserHandler{
		UserService: mockSvc,
	}

	router := gin.Default()
	router.POST("/login", userHandler.Login)

	// --- success ---
	mockSvc.EXPECT().
		AuthenticateUser(gomock.Any(), "test@example.com", "password123").
		Return("mocked.jwt.token", nil)

	loginReq := handler.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(loginReq)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set(constant.HeaderContentType, constant.ContentTypeJSON)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "mocked.jwt.token")

	// --- auth failed ---
	mockSvc.EXPECT().
		AuthenticateUser(gomock.Any(), "fail@example.com", "wrongpass").
		Return("", service.ErrInvalidCredential)

	failReq := handler.LoginRequest{
		Email:    "fail@example.com",
		Password: "wrongpass",
	}
	bodyFail, _ := json.Marshal(failReq)

	req, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyFail))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), service.ErrInvalidCredential.Error())

	// --- invalid payload ---
	req, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(`invalid`)))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
