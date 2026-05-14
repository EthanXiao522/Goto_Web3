package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xyd/web3-learning-tracker/internal/repository"
	"github.com/xyd/web3-learning-tracker/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	userRepo    *repository.UserRepo
}

func NewAuthHandler(authService *service.AuthService, userRepo *repository.UserRepo) *AuthHandler {
	return &AuthHandler{authService: authService, userRepo: userRepo}
}

type registerReq struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid input: " + err.Error()})
		return
	}
	user, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrDuplicateEmail):
			c.JSON(http.StatusConflict, gin.H{"code": 409, "msg": "email already registered"})
		case errors.Is(err, service.ErrDuplicateUsername):
			c.JSON(http.StatusConflict, gin.H{"code": 409, "msg": "username already taken"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": 201, "msg": "ok", "data": gin.H{"user": user}})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid input: " + err.Error()})
		return
	}
	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCreds) {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": gin.H{"token": token, "user": user}})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint64("user_id")
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": gin.H{"user": user}})
}

func (h *AuthHandler) CheckUsername(c *gin.Context) {
	username := c.Query("username")
	if username == "" || len(username) < 2 {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"available": false, "reason": "用户名至少2个字符"}})
		return
	}
	_, err := h.userRepo.FindByUsername(username)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"available": false, "reason": "用户名已被占用"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"available": true}})
}
