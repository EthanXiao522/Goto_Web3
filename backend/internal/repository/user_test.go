package repository

import (
	"testing"

	"github.com/xyd/web3-learning-tracker/internal/config"
	"github.com/xyd/web3-learning-tracker/internal/database"
	"github.com/xyd/web3-learning-tracker/internal/model"
)

func setupRepo(t *testing.T) {
	t.Helper()
	cfg := config.Load()
	if err := database.Connect(cfg.DBDSN); err != nil {
		t.Fatalf("connect: %v", err)
	}
	if err := database.Migrate(); err != nil {
		t.Fatalf("migrate: %v", err)
	}
}

func TestUserRepo_CreateAndFind(t *testing.T) {
	setupRepo(t)
	defer database.Close()

	repo := &UserRepo{DB: database.DB}
	name := "test_create_find"
	email := name + "@test.com"
	u := &model.User{Username: name, Email: email, PasswordHash: "hash"}
	id, err := repo.Create(u)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero id")
	}

	found, err := repo.FindByID(id)
	if err != nil {
		t.Fatalf("find by id: %v", err)
	}
	if found.Email != email {
		t.Errorf("email mismatch: %s", found.Email)
	}

	_, err = repo.FindByEmail(email)
	if err != nil {
		t.Fatalf("find by email: %v", err)
	}

	_, err = repo.FindByEmail("nonexist@test.com")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestPhaseRepo_GetAll(t *testing.T) {
	setupRepo(t)
	defer database.Close()

	repo := &PhaseRepo{DB: database.DB}
	phases, err := repo.GetAllWithProgress(1)
	if err != nil {
		t.Fatalf("get all: %v", err)
	}
	if len(phases) != 3 {
		t.Errorf("expected 3 phases, got %d", len(phases))
	}
	for _, p := range phases {
		if p.TaskCount == 0 {
			t.Errorf("Phase %d has 0 tasks", p.PhaseNumber)
		}
	}
}

func TestPhaseRepo_FindByID(t *testing.T) {
	setupRepo(t)
	defer database.Close()

	repo := &PhaseRepo{DB: database.DB}
	p, err := repo.FindByID(1)
	if err != nil {
		t.Fatalf("find by id: %v", err)
	}
	if p.PhaseNumber != 1 {
		t.Errorf("expected phase 1, got %d", p.PhaseNumber)
	}
}

func TestWeekRepo_FindByPhase(t *testing.T) {
	setupRepo(t)
	defer database.Close()

	repo := &WeekRepo{DB: database.DB}
	weeks, err := repo.FindByPhase(1)
	if err != nil {
		t.Fatalf("find by phase: %v", err)
	}
	if len(weeks) != 4 {
		t.Errorf("expected 4 weeks, got %d", len(weeks))
	}
}

func TestDayRepo_FindByWeek(t *testing.T) {
	setupRepo(t)
	defer database.Close()

	repo := &DayRepo{DB: database.DB}
	days, err := repo.FindByWeek(1)
	if err != nil {
		t.Fatalf("find by week: %v", err)
	}
	if len(days) != 7 {
		t.Errorf("expected 7 days, got %d", len(days))
	}
}

func TestTaskRepo_FindByDay(t *testing.T) {
	setupRepo(t)
	defer database.Close()

	repo := &TaskRepo{DB: database.DB}
	tasks, err := repo.FindByDay(1)
	if err != nil {
		t.Fatalf("find by day: %v", err)
	}
	if len(tasks) == 0 {
		t.Error("expected tasks in day 1")
	}
}

func TestUserTaskRepo_LazyCreate(t *testing.T) {
	setupRepo(t)
	defer database.Close()

	repo := &UserTaskRepo{DB: database.DB}
	ut, err := repo.LazyCreate(1, 1)
	if err != nil {
		t.Fatalf("lazy create: %v", err)
	}
	if ut.TaskID != 1 || ut.UserID != 1 {
		t.Errorf("unexpected user_task: %+v", ut)
	}

	ut2, err := repo.LazyCreate(1, 1)
	if err != nil {
		t.Fatalf("second lazy create: %v", err)
	}
	if ut2.ID != ut.ID {
		t.Error("lazy create not idempotent")
	}
}

func TestUserTaskRepo_UpdateComplete(t *testing.T) {
	setupRepo(t)
	defer database.Close()

	repo := &UserTaskRepo{DB: database.DB}
	if err := repo.UpdateComplete(1, 1, true); err != nil {
		t.Fatalf("update complete: %v", err)
	}

	ut, err := repo.FindByUserAndTask(1, 1)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if !ut.IsCompleted {
		t.Error("expected completed")
	}

	if err := repo.UpdateComplete(1, 1, false); err != nil {
		t.Fatalf("update uncomplete: %v", err)
	}
}
