package handler

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	gomarkdown "github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
	mdparser "github.com/gomarkdown/markdown/parser"

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
	jwtSecret       string
}

func NewPageHandler(
	userRepo *repository.UserRepo, phaseRepo *repository.PhaseRepo,
	weekRepo *repository.WeekRepo, dayRepo *repository.DayRepo,
	taskRepo *repository.TaskRepo, userTaskRepo *repository.UserTaskRepo,
	progressService *service.ProgressService, jwtSecret string,
) *PageHandler {
	return &PageHandler{
		userRepo: userRepo, phaseRepo: phaseRepo, weekRepo: weekRepo,
		dayRepo: dayRepo, taskRepo: taskRepo, userTaskRepo: userTaskRepo,
		progressService: progressService, jwtSecret: jwtSecret,
	}
}

func (h *PageHandler) Landing(c *gin.Context) {
	c.HTML(http.StatusOK, "landing.html", gin.H{})
}

func (h *PageHandler) LoginPage(c *gin.Context) {
	if h.isLoggedIn(c) {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}
	c.HTML(http.StatusOK, "auth_login.html", gin.H{})
}

func (h *PageHandler) RegisterPage(c *gin.Context) {
	if h.isLoggedIn(c) {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}
	c.HTML(http.StatusOK, "auth_register.html", gin.H{})
}

func (h *PageHandler) Dashboard(c *gin.Context) {
	userID := c.GetUint64("user_id")
	data, _ := h.progressService.GetDashboard(userID)
	phases, _ := h.phaseRepo.GetAllWithProgress(userID)

	// Gantt: weeks with progress
	var ganttPhases []gin.H
	for _, p := range phases {
		weeks, _ := h.weekRepo.FindByPhaseWithProgress(p.ID, userID)
		var wi []gin.H
		for _, w := range weeks {
			wi = append(wi, gin.H{
				"WeekNumber": w.WeekNumber, "Title": w.Title,
				"TaskCount": w.TaskCount, "CompletedCount": w.CompletedCount,
			})
		}
		ganttPhases = append(ganttPhases, gin.H{
			"ID": p.ID, "PhaseNumber": p.PhaseNumber, "Title": p.Title,
			"WeekCount": len(weeks), "Weeks": wi,
			"TaskCount": p.TaskCount, "CompletedCount": p.CompletedCount,
		})
	}

	// Day gantt: group days under weeks
	var dayGanttPhases []gin.H
	for _, p := range phases {
		weeks, _ := h.weekRepo.FindByPhase(p.ID)
		var dayWeeks []gin.H
		for _, w := range weeks {
			days, _ := h.dayRepo.FindByWeekWithProgress(w.ID, userID)
			var di []gin.H
			for _, d := range days {
				di = append(di, gin.H{
					"DayNumber": d.DayNumber, "Title": d.Title,
					"TaskCount": d.TaskCount, "CompletedCount": d.CompletedCount,
				})
			}
			dayWeeks = append(dayWeeks, gin.H{
				"WeekNumber": w.WeekNumber, "Days": di, "DayCount": len(days),
			})
		}
		dayGanttPhases = append(dayGanttPhases, gin.H{
			"ID": p.ID, "PhaseNumber": p.PhaseNumber, "Title": p.Title,
			"Weeks": dayWeeks,
		})
	}

	c.HTML(http.StatusOK, "dashboard.html", h.baseData(c, "Dashboard", "dashboard", gin.H{
		"Data":          data,
		"Phases":        phases,
		"GanttPhases":   ganttPhases,
		"DayGanttPhases": dayGanttPhases,
	}))
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

	var enriched []gin.H
	for _, p := range phases {
		weeks, _ := h.weekRepo.FindByPhaseWithProgress(p.ID, userID)
		var wi []gin.H
		for _, w := range weeks {
			wi = append(wi, gin.H{
				"WeekNumber": w.WeekNumber, "Title": w.Title,
				"TaskCount": w.TaskCount, "CompletedCount": w.CompletedCount,
			})
		}
		p.WeekCount = len(weeks)
		enriched = append(enriched, gin.H{
			"ID": p.ID, "PhaseNumber": p.PhaseNumber, "Title": p.Title,
			"SortOrder": p.SortOrder, "WeekCount": len(weeks), "Weeks": wi,
			"TaskCount": p.TaskCount, "CompletedCount": p.CompletedCount,
		})
	}

	c.HTML(http.StatusOK, "gantt.html", h.baseData(c, "甘特图", "gantt", gin.H{"Phases": enriched}))
}

func (h *PageHandler) LearningTasks(c *gin.Context) {
	userID := c.GetUint64("user_id")
	phases, _ := h.phaseRepo.GetAllWithProgress(userID)

	var enrichedPhases []gin.H
	for _, p := range phases {
		weeks, _ := h.weekRepo.FindByPhase(p.ID)
		var enrichedWeeks []gin.H
		var allTaskIDs []uint64
		for _, w := range weeks {
			days, _ := h.dayRepo.FindByWeek(w.ID)
			var enrichedDays []gin.H
			for _, d := range days {
				tasks, _ := h.taskRepo.FindByDay(d.ID)
				var enrichedTasks []gin.H
				for _, t := range tasks {
					allTaskIDs = append(allTaskIDs, t.ID)
					enrichedTasks = append(enrichedTasks, gin.H{
						"ID": t.ID, "Content": t.Content, "EstimatedHours": t.EstimatedHours,
						"ResourceURLs": t.ResourceURLs, "IsCheckpoint": t.IsCheckpoint, "SortOrder": t.SortOrder,
					})
				}
				enrichedDays = append(enrichedDays, gin.H{
					"ID": d.ID, "DayNumber": d.DayNumber, "Title": d.Title,
					"Tasks": enrichedTasks,
				})
			}
			enrichedWeeks = append(enrichedWeeks, gin.H{
				"ID": w.ID, "WeekNumber": w.WeekNumber, "Title": w.Title,
				"Days": enrichedDays,
			})
		}

		utMap, _ := h.userTaskRepo.FindByUserAndTaskIDs(userID, allTaskIDs)
		// Embed IsCompleted into each task
		for _, wk := range enrichedWeeks {
			for _, dy := range wk["Days"].([]gin.H) {
				for _, tk := range dy["Tasks"].([]gin.H) {
					taskID := tk["ID"].(uint64)
					if ut, ok := utMap[taskID]; ok {
						tk["IsCompleted"] = ut.IsCompleted
					}
				}
			}
		}

		enrichedPhases = append(enrichedPhases, gin.H{
			"ID": p.ID, "PhaseNumber": p.PhaseNumber, "Title": p.Title,
			"TaskCount": p.TaskCount, "CompletedCount": p.CompletedCount,
			"Weeks": enrichedWeeks, "UserTasks": utMap,
		})
	}

	c.HTML(http.StatusOK, "learning_tasks.html", h.baseData(c, "学习任务", "tasks", gin.H{
		"Phases": enrichedPhases,
	}))
}

var reTOC = regexp.MustCompile(`(?m)^\[TOC\]\s*\n?`)

func (h *PageHandler) Handbook(c *gin.Context) {
	data, err := os.ReadFile("../sources/web3_infra_3month_plan.md")
	if err != nil {
		slog.Error("read handbook md failed", "err", err)
		c.HTML(http.StatusOK, "handbook.html", h.baseData(c, "学习计划书", "handbook", gin.H{
			"Content": template.HTML("<p>手册内容加载失败</p>"),
			"TOC":     []gin.H{},
		}))
		return
	}
	raw := string(data)

	// Strip [TOC] placeholder from markdown
	clean := reTOC.ReplaceAllString(raw, "")

	extensions := mdparser.CommonExtensions | mdparser.AutoHeadingIDs
	parser := mdparser.NewWithExtensions(extensions)
	doc := parser.Parse([]byte(clean))
	renderer := mdhtml.NewRenderer(mdhtml.RendererOptions{
		Flags: mdhtml.CommonFlags,
	})
	html := gomarkdown.Render(doc, renderer)

	c.HTML(http.StatusOK, "handbook.html", h.baseData(c, "学习计划书", "handbook", gin.H{
		"Content": template.HTML(html),
	}))
}

func (h *PageHandler) HandbookSource(c *gin.Context) {
	c.File("../sources/web3_infra_3month_plan.md")
}

func (h *PageHandler) Demo(c *gin.Context) {
	phases, _ := h.phaseRepo.GetAllWithProgress(0)
	c.HTML(http.StatusOK, "demo.html", gin.H{
		"Title":   "游客模式",
		"Phases":  phases,
		"MockData": mockDashboardData(),
	})
}

func (h *PageHandler) ProfilePage(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.html", h.baseData(c, "修改个人信息", "profile", nil))
}

func mockDashboardData() gin.H {
	return gin.H{
		"overview": gin.H{
			"total_tasks":     215,
			"completed_tasks": 87,
			"total_phases":    3,
			"completed_phases": 1,
			"total_weeks":     12,
			"completed_weeks": 3,
		},
		"phase_progress": []gin.H{
			{"phase_id": 1, "phase_number": 1, "title": "基础夯实（第 1-4 周）", "task_count": 74, "completed_count": 52, "percentage": 70.0},
			{"phase_id": 2, "phase_number": 2, "title": "核心系统（第 5-8 周）", "task_count": 69, "completed_count": 28, "percentage": 40.0},
			{"phase_id": 3, "phase_number": 3, "title": "工程化与面试（第 9-12 周）", "task_count": 72, "completed_count": 7, "percentage": 10.0},
		},
		"week_progress": []gin.H{
			{"week_number": 1, "title": "Go 并发 + EVM 基础", "task_count": 18, "completed_count": 18},
			{"week_number": 2, "title": "Go 网络编程 + 区块扫描器", "task_count": 19, "completed_count": 16},
			{"week_number": 3, "title": "Listener 升级", "task_count": 18, "completed_count": 12},
			{"week_number": 4, "title": "Listener 打磨", "task_count": 19, "completed_count": 6},
			{"week_number": 5, "title": "Kafka 深化", "task_count": 18, "completed_count": 0},
		},
		"recent_tasks": []gin.H{
			{"content": "实现 WebSocket 自动重连 + 指数退避", "phase_title": "阶段 1 基础夯实", "week_number": 3, "is_completed": true, "completed_at": "2026-03-21 16:30"},
			{"content": "创建 event_logs 表并解析 ERC-20 Transfer", "phase_title": "阶段 1 基础夯实", "week_number": 2, "is_completed": true, "completed_at": "2026-03-14 15:20"},
			{"content": "实现 goroutine + channel 生产者-消费者模型", "phase_title": "阶段 1 基础夯实", "week_number": 1, "is_completed": true, "completed_at": "2026-03-07 11:45"},
			{"content": "Kafka producer 集成到 Listener", "phase_title": "阶段 2 核心系统", "week_number": 5, "is_completed": false},
			{"content": "设计 HD 钱包密钥管理方案", "phase_title": "阶段 2 核心系统", "week_number": 6, "is_completed": false},
		},
	}
}

func (h *PageHandler) baseData(c *gin.Context, title, active string, extra gin.H) gin.H {
	userID := c.GetUint64("user_id")
	user, _ := h.userRepo.FindByID(userID)
	phases, _ := h.phaseRepo.GetAllWithProgress(userID)

	var sidebarPhases []gin.H
	for _, p := range phases {
		sidebarPhases = append(sidebarPhases, gin.H{
			"ID": p.ID, "Title": "阶段 " + strconv.Itoa(int(p.PhaseNumber)) + " " + p.Title,
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

func (h *PageHandler) isLoggedIn(c *gin.Context) bool {
	tokenStr, err := c.Cookie("token")
	if err != nil {
		return false
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})
	return err == nil && token.Valid
}
