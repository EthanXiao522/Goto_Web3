package service

import (
	"fmt"

	"github.com/xyd/web3-learning-tracker/internal/model"
	"github.com/xyd/web3-learning-tracker/internal/repository"
)

type TaskService struct {
	taskRepo     *repository.TaskRepo
	userTaskRepo *repository.UserTaskRepo
}

func NewTaskService(taskRepo *repository.TaskRepo, userTaskRepo *repository.UserTaskRepo) *TaskService {
	return &TaskService{taskRepo: taskRepo, userTaskRepo: userTaskRepo}
}

func (s *TaskService) ToggleComplete(userID, taskID uint64, completed bool) (*model.UserTask, error) {
	ut, err := s.getOrCreateUserTask(userID, taskID)
	if err != nil {
		return nil, err
	}
	if err := s.userTaskRepo.UpdateComplete(userID, taskID, completed); err != nil {
		return nil, fmt.Errorf("task: toggle complete: %w", err)
	}
	ut.IsCompleted = completed
	return ut, nil
}

func (s *TaskService) UpdateSubmissions(userID, taskID uint64, fields map[string]string) (*model.UserTask, error) {
	ut, err := s.getOrCreateUserTask(userID, taskID)
	if err != nil {
		return nil, err
	}
	if err := s.userTaskRepo.UpdateFields(userID, taskID, fields); err != nil {
		return nil, fmt.Errorf("task: update submissions: %w", err)
	}
	return ut, nil
}

func (s *TaskService) getOrCreateUserTask(userID, taskID uint64) (*model.UserTask, error) {
	return s.userTaskRepo.LazyCreate(userID, taskID)
}
