package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/xyd/web3-learning-tracker/internal/config"
	"github.com/xyd/web3-learning-tracker/internal/database"
	"github.com/xyd/web3-learning-tracker/internal/model"
	"github.com/xyd/web3-learning-tracker/internal/repository"
	"github.com/xyd/web3-learning-tracker/internal/service"
)

func setupService(t *testing.T) uint64 {
	t.Helper()
	cfg := config.Load()
	if err := database.Connect(cfg.DBDSN); err != nil {
		t.Fatalf("connect: %v", err)
	}
	if err := database.Migrate(); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	userRepo := &repository.UserRepo{DB: database.DB}
	u := &model.User{Username: fmt.Sprintf("svcusr_%d", time.Now().UnixNano()), Email: fmt.Sprintf("svc_%d@test.com", time.Now().UnixNano()), PasswordHash: "$2a$12$dummy"}
	id, err := userRepo.Create(u)
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}
	return id
}

func TestProgressService_GetDashboard(t *testing.T) {
	uid := setupService(t)
	defer database.Close()

	svc := service.NewProgressService(database.DB)
	data, err := svc.GetDashboard(uid)
	if err != nil {
		t.Fatalf("get dashboard: %v", err)
	}
	if data.Overview.TotalTasks == 0 {
		t.Error("expected non-zero total tasks")
	}
	if data.Overview.TotalPhases != 3 {
		t.Errorf("expected 3 phases, got %d", data.Overview.TotalPhases)
	}
	if len(data.WeekProgress) == 0 {
		t.Error("expected non-empty week progress")
	}
}

func TestProgressService_GetOverview(t *testing.T) {
	uid := setupService(t)
	defer database.Close()

	svc := service.NewProgressService(database.DB)
	overview, err := svc.GetOverview(uid)
	if err != nil {
		t.Fatalf("get overview: %v", err)
	}
	if overview.TotalTasks == 0 {
		t.Error("expected non-zero total tasks")
	}
	if overview.TotalPhases != 3 {
		t.Errorf("expected 3 total phases, got %d", overview.TotalPhases)
	}
}

func TestTaskService_ToggleComplete(t *testing.T) {
	uid := setupService(t)
	defer database.Close()

	taskRepo := &repository.TaskRepo{DB: database.DB}
	userTaskRepo := &repository.UserTaskRepo{DB: database.DB}
	svc := service.NewTaskService(taskRepo, userTaskRepo)

	ut, err := svc.ToggleComplete(uid, 1, true)
	if err != nil {
		t.Fatalf("toggle complete: %v", err)
	}
	if !ut.IsCompleted {
		t.Error("expected completed")
	}

	ut, err = svc.ToggleComplete(uid, 1, false)
	if err != nil {
		t.Fatalf("toggle uncomplete: %v", err)
	}
	if ut.IsCompleted {
		t.Error("expected uncompleted")
	}
}

func TestTaskService_UpdateSubmissions(t *testing.T) {
	uid := setupService(t)
	defer database.Close()

	taskRepo := &repository.TaskRepo{DB: database.DB}
	userTaskRepo := &repository.UserTaskRepo{DB: database.DB}
	svc := service.NewTaskService(taskRepo, userTaskRepo)

	fields := map[string]string{
		"learning_links":      "https://svc.example.com",
		"implementation_plan": "svc test plan",
		"implementation_code": "svc test code",
		"experience_summary":  "svc test summary",
	}
	_, err := svc.UpdateSubmissions(uid, 1, fields)
	if err != nil {
		t.Fatalf("update submissions: %v", err)
	}

	// Re-fetch to verify persistence
	ut, err := userTaskRepo.FindByUserAndTask(uid, 1)
	if err != nil {
		t.Fatalf("find after update: %v", err)
	}
	if ut.LearningLinks != "https://svc.example.com" {
		t.Errorf("expected learning_links, got: %s", ut.LearningLinks)
	}
	if ut.ImplementationPlan != "svc test plan" {
		t.Errorf("expected implementation_plan, got: %s", ut.ImplementationPlan)
	}
}
