package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/xyd/web3-learning-tracker/internal/config"
	"github.com/xyd/web3-learning-tracker/internal/database"
	"github.com/xyd/web3-learning-tracker/internal/repository"
	"github.com/xyd/web3-learning-tracker/internal/service"
)

func setupHandler(t *testing.T) (*gin.Engine, string) {
	t.Helper()
	cfg := config.Load()
	if err := database.Connect(cfg.DBDSN); err != nil {
		t.Fatalf("connect: %v", err)
	}
	if err := database.Migrate(); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() { database.Close() })

	gin.SetMode(gin.TestMode)
	r := gin.New()

	userRepo := &repository.UserRepo{DB: database.DB}
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := NewAuthHandler(authService, userRepo)

	r.POST("/api/v1/auth/register", authHandler.Register)
	r.POST("/api/v1/auth/login", authHandler.Login)
	r.GET("/api/v1/auth/me", func(c *gin.Context) {
		c.Set("user_id", uint64(1))
		c.Next()
	}, authHandler.Me)

	return r, cfg.JWTSecret
}

func TestAuthHandler_Register(t *testing.T) {
	r, _ := setupHandler(t)

	body := map[string]string{"username": "handler_reg_test", "email": "handler_reg@test.com", "password": "123456"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) != 201 {
		t.Errorf("expected code 201, got %v", resp["code"])
	}
}

func TestAuthHandler_RegisterDuplicate(t *testing.T) {
	r, _ := setupHandler(t)

	body := map[string]string{"username": "handler_dup_test", "email": "handler_dup@test.com", "password": "123456"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(b))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)

	if w2.Code != 409 {
		t.Errorf("expected 409, got %d", w2.Code)
	}
}

func TestAuthHandler_Login(t *testing.T) {
	r, _ := setupHandler(t)

	body := map[string]string{"email": "test@test.com", "password": "123456"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["token"].(string) == "" {
		t.Error("expected token")
	}
}

func TestAuthHandler_LoginInvalid(t *testing.T) {
	r, _ := setupHandler(t)

	body := map[string]string{"email": "test@test.com", "password": "wrongpass"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthHandler_Me(t *testing.T) {
	r, _ := setupHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
