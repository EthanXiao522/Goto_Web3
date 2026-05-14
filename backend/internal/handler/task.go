package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/xyd/web3-learning-tracker/internal/repository"
	"github.com/xyd/web3-learning-tracker/internal/service"
)

type TaskHandler struct {
	taskService *service.TaskService
	taskRepo    *repository.TaskRepo
	userTaskRepo *repository.UserTaskRepo
}

func NewTaskHandler(taskService *service.TaskService, taskRepo *repository.TaskRepo, userTaskRepo *repository.UserTaskRepo) *TaskHandler {
	return &TaskHandler{taskService: taskService, taskRepo: taskRepo, userTaskRepo: userTaskRepo}
}

type toggleReq struct {
	Completed bool `json:"completed"`
}

func (h *TaskHandler) ToggleComplete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid task id"})
		return
	}
	var req toggleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid input"})
		return
	}
	ut, err := h.taskService.ToggleComplete(userID, taskID, req.Completed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user_task": ut}})
}

type submissionsReq struct {
	LearningLinks     string `json:"learning_links"`
	ImplementationPlan string `json:"implementation_plan"`
	ImplementationCode string `json:"implementation_code"`
	ExperienceSummary  string `json:"experience_summary"`
}

func (h *TaskHandler) UpdateSubmissions(c *gin.Context) {
	userID := c.GetUint64("user_id")
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid task id"})
		return
	}
	var req submissionsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid input"})
		return
	}
	fields := map[string]string{
		"learning_links":      req.LearningLinks,
		"implementation_plan": req.ImplementationPlan,
		"implementation_code": req.ImplementationCode,
		"experience_summary":  req.ExperienceSummary,
	}
	ut, err := h.taskService.UpdateSubmissions(userID, taskID, fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user_task": ut}})
}

func (h *TaskHandler) GetTaskDetail(c *gin.Context) {
	userID := c.GetUint64("user_id")
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid task id"})
		return
	}
	task, err := h.taskRepo.FindByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "task not found"})
		return
	}
	ut, _ := h.userTaskRepo.FindByUserAndTask(userID, taskID)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"task": task, "user_task": ut}})
}
