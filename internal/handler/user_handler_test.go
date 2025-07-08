package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SoliMark/gotasker-pro/internal/handler"
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
