package router

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"math"

	"github.com/gin-gonic/gin"

	"github.com/xyd/web3-learning-tracker/internal/handler"
	"github.com/xyd/web3-learning-tracker/internal/middleware"
	"github.com/xyd/web3-learning-tracker/internal/model"
	"github.com/xyd/web3-learning-tracker/internal/repository"
	"github.com/xyd/web3-learning-tracker/internal/service"
)

func Setup(db *sql.DB, jwtSecret string) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Template functions
	funcMap := template.FuncMap{
		"percent": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return int(math.Round(float64(a) / float64(b) * 100))
		},
		"multiply": func(a float64, b float64) float64 { return a * b },
		"phaseColor": func(num any) string {
			n := toInt(num)
			switch n {
			case 1:
				return "#00f0ff"
			case 2:
				return "#b347ea"
			case 3:
				return "#00ff88"
			}
			return "#00f0ff"
		},
		"phaseColorClass": func(num any) string {
			n := toInt(num)
			switch n {
			case 1:
				return "cyan"
			case 2:
				return "purple"
			case 3:
				return "green"
			}
			return "cyan"
		},
		"resourceURLs": func(raw string) []string {
			var urls []string
			if raw == "" || raw == "[]" {
				return urls
			}
			json.Unmarshal([]byte(raw), &urls)
			return urls
		},
		"ganttLeft": func(phaseNum uint8) float64 {
			return float64(phaseNum-1) * 33.33
		},
		"ganttWidth": func(weekCount int) float64 {
			return float64(weekCount) / 12.0 * 100
		},
		"ganttWeekLeft": func(weekNum int) float64 {
			return float64(weekNum-1) / 12.0 * 100
		},
		"fieldValue": func(ut *model.UserTask, field string) string {
			if ut == nil {
				return ""
			}
			switch field {
			case "learning_links":
				return ut.LearningLinks
			case "implementation_plan":
				return ut.ImplementationPlan
			case "implementation_code":
				return ut.ImplementationCode
			case "experience_summary":
				return ut.ExperienceSummary
			}
			return ""
		},
	}
	r.SetFuncMap(funcMap)
	r.LoadHTMLGlob("../frontend/templates/*.html")
	r.Static("/static", "../frontend/static")

	// Repos
	userRepo := &repository.UserRepo{DB: db}
	phaseRepo := &repository.PhaseRepo{DB: db}
	weekRepo := &repository.WeekRepo{DB: db}
	dayRepo := &repository.DayRepo{DB: db}
	taskRepo := &repository.TaskRepo{DB: db}
	userTaskRepo := &repository.UserTaskRepo{DB: db}

	// Services
	authService := service.NewAuthService(userRepo, jwtSecret)
	taskService := service.NewTaskService(taskRepo, userTaskRepo)
	progressService := service.NewProgressService(db)

	// Handlers
	authHandler := handler.NewAuthHandler(authService, userRepo)
	phaseHandler := handler.NewPhaseHandler(phaseRepo, weekRepo, dayRepo)
	taskHandler := handler.NewTaskHandler(taskService, taskRepo, userTaskRepo)
	progressHandler := handler.NewProgressHandler(progressService)
	pageHandler := handler.NewPageHandler(userRepo, phaseRepo, weekRepo, dayRepo, taskRepo, userTaskRepo, progressService, jwtSecret)

	// Public page routes
	r.GET("/", pageHandler.Landing)
	r.GET("/login", pageHandler.LoginPage)
	r.GET("/register", pageHandler.RegisterPage)
	r.GET("/demo", pageHandler.Demo)
	r.GET("/logout", func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/", "", false, true)
		c.Redirect(302, "/")
	})

	// API routes
	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.GET("/auth/check-username", authHandler.CheckUsername)

		protected := api.Group("")
		protected.Use(middleware.Auth(jwtSecret))
		{
			protected.GET("/auth/me", authHandler.Me)
				protected.PUT("/auth/profile", authHandler.UpdateProfile)
			protected.GET("/phases", phaseHandler.GetPhases)
			protected.GET("/phases/:id", phaseHandler.GetPhaseDetail)
			protected.GET("/weeks/:id", phaseHandler.GetWeekDetail)
			protected.GET("/tasks/:id", taskHandler.GetTaskDetail)
			protected.PATCH("/tasks/:id/complete", taskHandler.ToggleComplete)
			protected.PUT("/tasks/:id/submissions", taskHandler.UpdateSubmissions)
			protected.GET("/dashboard", progressHandler.GetDashboard)
			protected.GET("/progress", progressHandler.GetOverview)
		}
	}

	// Authenticated page routes
	authPages := r.Group("")
	authPages.Use(middleware.Auth(jwtSecret))
	{
		authPages.GET("/dashboard", pageHandler.Dashboard)
		authPages.GET("/phases", pageHandler.PhaseList)
		authPages.GET("/phases/:id", pageHandler.PhaseDetail)
		authPages.GET("/weeks/:id", pageHandler.WeekDetail)
		authPages.GET("/tasks/:id", pageHandler.TaskDetail)
		authPages.GET("/gantt", pageHandler.Gantt)
		authPages.GET("/handbook", pageHandler.Handbook)
			authPages.GET("/profile", pageHandler.ProfilePage)
	}

	return r
}

func toInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case uint8:
		return int(n)
	case uint64:
		return int(n)
	case float64:
		return int(n)
	}
	return 0
}

