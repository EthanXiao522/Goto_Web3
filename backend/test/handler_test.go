package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/xyd/web3-learning-tracker/internal/config"
	"github.com/xyd/web3-learning-tracker/internal/database"
	"github.com/xyd/web3-learning-tracker/internal/handler"
	"github.com/xyd/web3-learning-tracker/internal/repository"
	"github.com/xyd/web3-learning-tracker/internal/service"
)

func uniqueNameH(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

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
	authHandler := handler.NewAuthHandler(authService, userRepo)

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

	uname := uniqueNameH("handler_reg_test")
	body := map[string]string{"username": uname, "email": uname + "@test.com", "password": "123456"}
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

	uname := uniqueNameH("handler_dup_test")
	body := map[string]string{"username": uname, "email": uname + "@test.com", "password": "123456"}
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

func setupTaskHandler(t *testing.T) *gin.Engine {
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

	taskRepo := &repository.TaskRepo{DB: database.DB}
	userTaskRepo := &repository.UserTaskRepo{DB: database.DB}
	taskService := service.NewTaskService(taskRepo, userTaskRepo)
	taskHandler := handler.NewTaskHandler(taskService, taskRepo, userTaskRepo)

	auth := func(c *gin.Context) {
		c.Set("user_id", uint64(1))
		c.Next()
	}

	protected := r.Group("")
	protected.Use(auth)
	{
		protected.GET("/api/v1/tasks/:id", taskHandler.GetTaskDetail)
		protected.PATCH("/api/v1/tasks/:id/complete", taskHandler.ToggleComplete)
		protected.PUT("/api/v1/tasks/:id/content", taskHandler.UpdateContent)
		protected.PUT("/api/v1/tasks/:id/submissions", taskHandler.UpdateSubmissions)
	}

	return r
}

func TestTaskHandler_GetTaskDetail(t *testing.T) {
	r := setupTaskHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/tasks/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["task"] == nil {
		t.Error("expected task in response")
	}
}

func TestTaskHandler_ToggleComplete(t *testing.T) {
	r := setupTaskHandler(t)

	body := map[string]bool{"completed": true}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PATCH", "/api/v1/tasks/1/complete", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Toggle back
	body["completed"] = false
	b, _ = json.Marshal(body)
	req2 := httptest.NewRequest("PATCH", "/api/v1/tasks/1/complete", bytes.NewReader(b))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != 200 {
		t.Errorf("uncomplete: expected 200, got %d", w2.Code)
	}
}

func TestTaskHandler_UpdateContent(t *testing.T) {
	r := setupTaskHandler(t)

	body := map[string]string{"content": "[TEST] updated task content"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/api/v1/tasks/1/content", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify update
	req2 := httptest.NewRequest("GET", "/api/v1/tasks/1", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	var resp map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	task := resp["data"].(map[string]interface{})["task"].(map[string]interface{})
	if task["content"].(string) != "[TEST] updated task content" {
		t.Errorf("content not updated: %s", task["content"])
	}
}

func TestTaskHandler_UpdateContentEmpty(t *testing.T) {
	r := setupTaskHandler(t)

	body := map[string]string{"content": ""}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/api/v1/tasks/1/content", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400 for empty content, got %d", w.Code)
	}
}

func TestTaskHandler_UpdateSubmissions(t *testing.T) {
	r := setupTaskHandler(t)

	body := map[string]string{
		"learning_links":      "https://example.com",
		"implementation_plan": "plan text",
		"implementation_code": "code text",
		"experience_summary":  "summary text",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/api/v1/tasks/1/submissions", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func setupPhaseHandler(t *testing.T) *gin.Engine {
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

	phaseRepo := &repository.PhaseRepo{DB: database.DB}
	weekRepo := &repository.WeekRepo{DB: database.DB}
	dayRepo := &repository.DayRepo{DB: database.DB}
	phaseHandler := handler.NewPhaseHandler(phaseRepo, weekRepo, dayRepo)

	auth := func(c *gin.Context) {
		c.Set("user_id", uint64(1))
		c.Next()
	}

	protected := r.Group("")
	protected.Use(auth)
	{
		protected.GET("/api/v1/phases", phaseHandler.GetPhases)
		protected.GET("/api/v1/phases/:id", phaseHandler.GetPhaseDetail)
		protected.GET("/api/v1/weeks/:id", phaseHandler.GetWeekDetail)
	}

	return r
}

func TestPhaseHandler_GetPhases(t *testing.T) {
	r := setupPhaseHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/phases", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	phases := resp["data"].(map[string]interface{})["phases"].([]interface{})
	if len(phases) != 3 {
		t.Errorf("expected 3 phases, got %d", len(phases))
	}
}

func TestPhaseHandler_GetPhaseDetail(t *testing.T) {
	r := setupPhaseHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/phases/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	weeks := resp["data"].(map[string]interface{})["weeks"].([]interface{})
	if len(weeks) != 4 {
		t.Errorf("expected 4 weeks, got %d", len(weeks))
	}
}

func TestPhaseHandler_GetPhaseDetailNotFound(t *testing.T) {
	r := setupPhaseHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/phases/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestPhaseHandler_GetWeekDetail(t *testing.T) {
	r := setupPhaseHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/weeks/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	days := resp["data"].(map[string]interface{})["days"].([]interface{})
	if len(days) != 7 {
		t.Errorf("expected 7 days, got %d", len(days))
	}
}

func setupProgressHandler(t *testing.T) *gin.Engine {
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

	progressService := service.NewProgressService(database.DB)
	progressHandler := handler.NewProgressHandler(progressService)

	auth := func(c *gin.Context) {
		c.Set("user_id", uint64(1))
		c.Next()
	}

	protected := r.Group("")
	protected.Use(auth)
	{
		protected.GET("/api/v1/dashboard", progressHandler.GetDashboard)
		protected.GET("/api/v1/progress", progressHandler.GetOverview)
	}

	return r
}

func TestProgressHandler_GetDashboard(t *testing.T) {
	r := setupProgressHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/dashboard", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	overview := data["overview"].(map[string]interface{})
	if overview["total_tasks"].(float64) == 0 {
		t.Error("expected non-zero total_tasks")
	}
}

func TestProgressHandler_GetOverview(t *testing.T) {
	r := setupProgressHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/progress", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["data"] == nil {
		t.Error("expected data in response")
	}
}
