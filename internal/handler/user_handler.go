package handler

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/andrMaulana/Go-Task-Management-System/internal/models"
	"github.com/andrMaulana/Go-Task-Management-System/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Meta: models.Meta{
				Code:    http.StatusBadRequest,
				Status:  "Bad Request",
				Message: "Invalid input data",
			},
		})
		return
	}

	if err := validateUser(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Meta: models.Meta{
				Code:    http.StatusBadRequest,
				Status:  "Bad Request",
				Message: err.Error(),
			},
		})
		return
	}

	if err := h.userService.Register(&user); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Meta: models.Meta{
				Code:    http.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to register user",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Meta: models.Meta{
			Code:    http.StatusCreated,
			Status:  "Created",
			Message: "User registered successfully",
		},
		Data: gin.H{
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Meta: models.Meta{
				Code:    http.StatusBadRequest,
				Message: "Invalid input data",
			},
		})
		return
	}

	user, token, err := h.userService.Login(loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Response{
			Meta: models.Meta{
				Code:    http.StatusUnauthorized,
				Message: "Invalid credentials",
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Meta: models.Meta{
			Code:    http.StatusOK,
			Message: "Login successful",
		},
		Data: gin.H{
			"username": user.Username,
			"email":    user.Email,
			"token":    token,
		},
	})
}

func validateUser(user *models.User) error {
	if len(user.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(user.Email) {
		return errors.New("invalid email format")
	}

	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}
