package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SoliMark/gotasker-pro/internal/constant"
	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/service"
)

type UserHandler struct {
	UserService service.UserService
}

// ==================== Register ====================

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

// ==================== Login ====================

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request:" + err.Error(),
		})
		return
	}

	user := &model.User{
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	if err := h.UserService.CreateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	token, err := h.UserService.AuthenticateUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to generate token" + err.Error(),
		})
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request: " + err.Error(),
		})
		return
	}

	token, err := h.UserService.AuthenticateUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
	})
}

func (h *UserHandler) Profile(c *gin.Context) {
	userID, exists := c.Get(constant.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "user_id not found in context",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
	})
}
