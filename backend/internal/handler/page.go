package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/xyd/web3-learning-tracker/internal/model"
	"github.com/xyd/web3-learning-tracker/internal/repository"
	"github.com/xyd/web3-learning-tracker/internal/service"
)

type PageHandler struct {
	userRepo        *repository.UserRepo
	phaseRepo       *repository.PhaseRepo
	weekRepo        *repository.WeekRepo
	dayRepo         *repository.DayRepo
	taskRepo        *repository.TaskRepo
	userTaskRepo    *repository.UserTaskRepo
	progressService *service.ProgressService
}

func NewPageHandler(
	userRepo *repository.UserRepo, phaseRepo *repository.PhaseRepo,
	weekRepo *repository.WeekRepo, dayRepo *repository.DayRepo,
	taskRepo *repository.TaskRepo, userTaskRepo *repository.UserTaskRepo,
	progressService *service.ProgressService,
) *PageHandler {
	return &PageHandler{
		userRepo: userRepo, phaseRepo: phaseRepo, weekRepo: weekRepo,
		dayRepo: dayRepo, taskRepo: taskRepo, userTaskRepo: userTaskRepo,
		progressService: progressService,
	}
}

func (h *PageHandler) Landing(c *gin.Context) {
	c.HTML(http.StatusOK, "landing.html", nil)
}

func (h *PageHandler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth_login.html", nil)
}

func (h *PageHandler) RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth_register.html", nil)
}

func (h *PageHandler) Dashboard(c *gin.Context) {
	userID := c.GetUint64("user_id")
	data, _ := h.progressService.GetDashboard(userID)
	c.HTML(http.StatusOK, "dashboard.html", h.baseData(c, "Dashboard", "dashboard", gin.H{"Data": data}))
}

func (h *PageHandler) PhaseList(c *gin.Context) {
	userID := c.GetUint64("user_id")
	phases, _ := h.phaseRepo.GetAllWithProgress(userID)
	c.HTML(http.StatusOK, "phases.html", h.baseData(c, "学习阶段", "phases", gin.H{"Phases": phases}))
}

func (h *PageHandler) PhaseDetail(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	phase, _ := h.phaseRepo.FindByID(id)
	weeks, _ := h.weekRepo.FindByPhase(id)
	daysMap := make(map[uint64][]model.Day)
	tasksMap := make(map[uint64][]model.Task)
	utMap := make(map[uint64]*model.UserTask)

	var enrichedWeeks []gin.H
	for _, w := range weeks {
		days, _ := h.dayRepo.FindByWeek(w.ID)
		daysMap[w.ID] = days
		var taskIDs []uint64
		for _, d := range days {
			tasks, _ := h.taskRepo.FindByDay(d.ID)
			tasksMap[d.ID] = tasks
			for _, t := range tasks {
				taskIDs = append(taskIDs, t.ID)
			}
		}
		if len(taskIDs) > 0 {
			utMap, _ = h.userTaskRepo.FindByUserAndTaskIDs(userID, taskIDs)
		}

		var enrichedDays []gin.H
		for _, d := range days {
			var enrichedTasks []gin.H
			dayTaskCount := 0
			dayCompleted := 0
			for _, t := range tasksMap[d.ID] {
				ut := utMap[t.ID]
				completed := ut != nil && ut.IsCompleted
				hasSub := ut != nil && (ut.LearningLinks != "" || ut.ImplementationPlan != "" || ut.ImplementationCode != "" || ut.ExperienceSummary != "")
				if !t.IsCheckpoint {
					dayTaskCount++
					if completed {
						dayCompleted++
					}
				}
				enrichedTasks = append(enrichedTasks, gin.H{
					"ID": t.ID, "Content": t.Content, "IsCheckpoint": t.IsCheckpoint,
					"ResourceURLs": t.ResourceURLs, "UserCompleted": completed,
					"HasSubmission": hasSub,
				})
			}
			d.CompletedCount = dayCompleted
			d.TaskCount = dayTaskCount
			enrichedDays = append(enrichedDays, gin.H{
				"ID": d.ID, "WeekID": d.WeekID, "DayNumber": d.DayNumber,
				"Title": d.Title, "SortOrder": d.SortOrder,
				"TaskCount": dayTaskCount, "CompletedCount": dayCompleted,
				"Tasks": enrichedTasks,
			})
		}
		enrichedWeeks = append(enrichedWeeks, gin.H{
			"ID": w.ID, "PhaseID": w.PhaseID, "WeekNumber": w.WeekNumber,
			"Title": w.Title, "Subtitle": w.Subtitle, "Goal": w.Goal,
			"Deliverables": w.Deliverables, "SortOrder": w.SortOrder,
			"Days": enrichedDays,
		})
	}

	c.HTML(http.StatusOK, "phase_detail.html", h.baseData(c, phase.Title, "phases", gin.H{
		"Phase": phase, "Weeks": enrichedWeeks,
	}))
}

func (h *PageHandler) WeekDetail(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	week, _ := h.weekRepo.FindByID(id)
	days, _ := h.dayRepo.FindByWeek(id)
	phase, _ := h.phaseRepo.FindByID(week.PhaseID)

	var allTaskIDs []uint64
	for _, d := range days {
		tasks, _ := h.taskRepo.FindByDay(d.ID)
		for _, t := range tasks {
			allTaskIDs = append(allTaskIDs, t.ID)
		}
	}
	utMap, _ := h.userTaskRepo.FindByUserAndTaskIDs(userID, allTaskIDs)

	var enrichedDays []gin.H
	for _, d := range days {
		tasks, _ := h.taskRepo.FindByDay(d.ID)
		var enrichedTasks []gin.H
		dtCount, dtDone := 0, 0
		for _, t := range tasks {
			if t.IsCheckpoint {
				continue
			}
			dtCount++
			ut := utMap[t.ID]
			completed := ut != nil && ut.IsCompleted
			if completed {
				dtDone++
			}
			enrichedTasks = append(enrichedTasks, gin.H{
				"ID": t.ID, "Content": t.Content, "IsCheckpoint": t.IsCheckpoint,
				"ResourceURLs": t.ResourceURLs, "UserCompleted": completed,
			})
		}
		d.TaskCount = dtCount
		d.CompletedCount = dtDone
		enrichedDays = append(enrichedDays, gin.H{
			"ID": d.ID, "DayNumber": d.DayNumber, "Title": d.Title,
			"TaskCount": dtCount, "CompletedCount": dtDone, "Tasks": enrichedTasks,
		})
	}

	c.HTML(http.StatusOK, "week_detail.html", h.baseData(c, "第"+strconv.Itoa(int(week.WeekNumber))+"周", "phases", gin.H{
		"Week": week, "Phase": phase, "Days": enrichedDays,
	}))
}

func (h *PageHandler) TaskDetail(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	task, _ := h.taskRepo.FindByID(id)
	day, _ := h.dayRepo.FindByID(task.DayID)
	week, _ := h.weekRepo.FindByID(day.WeekID)
	ut, _ := h.userTaskRepo.FindByUserAndTask(userID, id)

	// Find prev/next task
	tasks, _ := h.taskRepo.FindByDay(task.DayID)
	var prevID, nextID uint64
	for i, t := range tasks {
		if t.ID == id {
			if i > 0 {
				prevID = tasks[i-1].ID
			}
			if i < len(tasks)-1 {
				nextID = tasks[i+1].ID
			}
			break
		}
	}

	c.HTML(http.StatusOK, "task_detail.html", h.baseData(c, "任务详情", "phases", gin.H{
		"Task": task, "UserTask": ut, "Day": day, "Week": week,
		"PrevTaskID": prevID, "NextTaskID": nextID,
	}))
}

func (h *PageHandler) Gantt(c *gin.Context) {
	userID := c.GetUint64("user_id")
	phases, _ := h.phaseRepo.GetAllWithProgress(userID)

	type weekInfo struct {
		WeekNumber int    `json:"week_number"`
		Title      string `json:"title"`
	}
	type phaseGantt struct {
		model.Phase
		Weeks []weekInfo `json:"weeks"`
	}

	var enriched []gin.H
	for _, p := range phases {
		weeks, _ := h.weekRepo.FindByPhase(p.ID)
		var wi []gin.H
		for _, w := range weeks {
			wi = append(wi, gin.H{"WeekNumber": w.WeekNumber, "Title": w.Title})
		}
		p.WeekCount = len(weeks)
		enriched = append(enriched, gin.H{
			"ID": p.ID, "PhaseNumber": p.PhaseNumber, "Title": p.Title,
			"SortOrder": p.SortOrder, "WeekCount": len(weeks), "Weeks": wi,
		})
	}

	c.HTML(http.StatusOK, "gantt.html", h.baseData(c, "甘特图", "gantt", gin.H{"Phases": enriched}))
}

func (h *PageHandler) Handbook(c *gin.Context) {
	c.HTML(http.StatusOK, "handbook.html", h.baseData(c, "学习手册", "handbook", gin.H{
		"Content": "<p>手册内容将通过 importer 导入后渲染</p>",
		"TOC":     []gin.H{},
	}))
}

func (h *PageHandler) baseData(c *gin.Context, title, active string, extra gin.H) gin.H {
	userID := c.GetUint64("user_id")
	user, _ := h.userRepo.FindByID(userID)
	phases, _ := h.phaseRepo.GetAllWithProgress(userID)

	var sidebarPhases []gin.H
	for _, p := range phases {
		sidebarPhases = append(sidebarPhases, gin.H{
			"ID": p.ID, "Title": "Phase " + strconv.Itoa(int(p.PhaseNumber)) + " " + p.Title,
			"TaskCount": p.TaskCount, "CompletedCount": p.CompletedCount,
		})
	}

	data := gin.H{
		"Title":         title,
		"Active":        active,
		"User":          user,
		"Breadcrumb":    title,
		"SidebarPhases": sidebarPhases,
	}
	for k, v := range extra {
		data[k] = v
	}
	return data
}
