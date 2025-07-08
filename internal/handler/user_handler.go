package handler

import (
	"net/http"

	"github.com/SoliMark/gotasker-pro/internal/model"
	"github.com/SoliMark/gotasker-pro/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService service.UserService
}

// input format
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// reponse success
type RegisterResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invaild request:" + err.Error(),
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

	c.JSON(http.StatusOK, RegisterResponse{
		Message: "User registered successfully",
	})
}
