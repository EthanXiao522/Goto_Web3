package handler

import (
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

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

type updateContentReq struct {
	Content string `json:"content"`
}

func (h *TaskHandler) UpdateContent(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid task id"})
		return
	}
	var req updateContentReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid input"})
		return
	}

	task, err := h.taskRepo.FindByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "task not found"})
		return
	}
	oldContent := task.Content

	if err := h.taskRepo.UpdateContent(taskID, req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	syncTaskToMd(oldContent, req.Content)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok"})
}

var mdLinkRe = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)

func stripMDLinks(s string) string {
	return mdLinkRe.ReplaceAllString(s, "[$1]")
}

func syncTaskToMd(oldContent, newContent string) {
	mdPath := "../sources/web3_infra_3month_plan.md"
	data, err := os.ReadFile(mdPath)
	if err != nil {
		slog.Error("sync md: read failed", "err", err)
		return
	}
	content := string(data)

	// Try exact match first
	if strings.Contains(content, oldContent) {
		updated := strings.Replace(content, oldContent, newContent, 1)
		if err := os.WriteFile(mdPath, []byte(updated), 0644); err != nil {
			slog.Error("sync md: write failed", "err", err)
			return
		}
		slog.Info("sync md: task content synced to md file (exact)")
		return
	}

	// Fallback: find line by stripping markdown links, then replace the whole task text
	lines := strings.Split(content, "\n")
	key := strings.TrimSpace(oldContent)
	if len(key) > 40 {
		key = key[:40]
	}
	for i, line := range lines {
		stripped := stripMDLinks(line)
		if strings.Contains(stripped, oldContent) || strings.Contains(stripped, key) {
			// Replace the task text (after "- [ ] ") in the original line
			prefix := ""
			trimmed := line
			if strings.HasPrefix(line, "- [ ] ") {
				prefix = "- [ ] "
				trimmed = line[6:]
			} else if strings.HasPrefix(line, "- [x] ") {
				prefix = "- [x] "
				trimmed = line[6:]
			}
			newLine := prefix + strings.Replace(trimmed, oldContent, newContent, 1)
			if newLine == line {
				// oldContent not a direct substring – replace entire task text
				newLine = prefix + newContent
			}
			lines[i] = newLine
			if err := os.WriteFile(mdPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
				slog.Error("sync md: write failed", "err", err)
				return
			}
			slog.Info("sync md: task content synced to md file (line match)")
			return
		}
	}
	slog.Warn("sync md: old content not found in md file", "key", key[:min(40, len(key))])
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
