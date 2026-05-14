package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/xyd/web3-learning-tracker/internal/repository"
)

type PhaseHandler struct {
	phaseRepo *repository.PhaseRepo
	weekRepo  *repository.WeekRepo
	dayRepo   *repository.DayRepo
}

func NewPhaseHandler(phaseRepo *repository.PhaseRepo, weekRepo *repository.WeekRepo, dayRepo *repository.DayRepo) *PhaseHandler {
	return &PhaseHandler{phaseRepo: phaseRepo, weekRepo: weekRepo, dayRepo: dayRepo}
}

func (h *PhaseHandler) GetPhases(c *gin.Context) {
	userID := c.GetUint64("user_id")
	phases, err := h.phaseRepo.GetAllWithProgress(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"phases": phases}})
}

func (h *PhaseHandler) GetPhaseDetail(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid phase id"})
		return
	}
	phase, err := h.phaseRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "phase not found"})
		return
	}
	weeks, err := h.weekRepo.FindByPhaseWithProgress(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"phase": phase, "weeks": weeks}})
}

func (h *PhaseHandler) GetWeekDetail(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid week id"})
		return
	}
	week, err := h.weekRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "week not found"})
		return
	}
	days, err := h.dayRepo.FindByWeekWithProgress(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"week": week, "days": days}})
}
