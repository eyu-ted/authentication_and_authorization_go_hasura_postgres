package controller

import (
	// domain "blog/models"
	"blog/services"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/hasura/go-graphql-client"
	// "blog/config"
	// "log"
	// "context"
	// "github.com/hasura/go-graphql-client"
)

type UserController struct {
	UserUsecase services.AuthUsecase
}

func NewAuthController(useruscase services.AuthUsecase) *UserController {
	return &UserController{useruscase}
}

func (h *UserController) Signup(c *gin.Context) {

	var request struct {
		Input struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		} `json:"input"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 2. Call your usecase (unchanged business logic)
	accessToken, refreshToken, err := h.UserUsecase.Signup(
		request.Input.Name,
		request.Input.Email,
		request.Input.Password,
	)

	// 3. Handle responses
	if err != nil {
		// Hasura-compatible error format
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"extensions": gin.H{
				"code": "SIGNUP_FAILED",
			},
		})
		return
	}

	// 4. Return Hasura-compatible success response
	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *UserController) Login(c *gin.Context) {

	var request struct {
		Input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		} `json:"input"`
	}
	// var user domain.User
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := request.Input.Email
	password := request.Input.Password

	accessToken, refreshToken, err := h.UserUsecase.Login(email, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *UserController) Refresh(c *gin.Context) {
	var refresh struct {
		RefreshToken string `json:"refresh"`
	}
	if err := c.ShouldBindJSON(&refresh); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	refreshToken := refresh.RefreshToken

	accessToken, err := h.UserUsecase.RefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}
