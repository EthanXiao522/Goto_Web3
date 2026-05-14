package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xyd/web3-learning-tracker/internal/service"
)

type ProgressHandler struct {
	progressService *service.ProgressService
}

func NewProgressHandler(progressService *service.ProgressService) *ProgressHandler {
	return &ProgressHandler{progressService: progressService}
}

func (h *ProgressHandler) GetDashboard(c *gin.Context) {
	userID := c.GetUint64("user_id")
	data, err := h.progressService.GetDashboard(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": data})
}

func (h *ProgressHandler) GetOverview(c *gin.Context) {
	userID := c.GetUint64("user_id")
	overview, err := h.progressService.GetOverview(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": overview})
}
