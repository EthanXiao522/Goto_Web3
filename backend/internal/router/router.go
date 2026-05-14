package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/xyd/web3-learning-tracker/internal/handler"
	"github.com/xyd/web3-learning-tracker/internal/middleware"
	"github.com/xyd/web3-learning-tracker/internal/repository"
	"github.com/xyd/web3-learning-tracker/internal/service"
)

func Setup(db *sql.DB, jwtSecret string) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	userRepo := &repository.UserRepo{DB: db}
	phaseRepo := &repository.PhaseRepo{DB: db}
	weekRepo := &repository.WeekRepo{DB: db}
	dayRepo := &repository.DayRepo{DB: db}
	taskRepo := &repository.TaskRepo{DB: db}
	userTaskRepo := &repository.UserTaskRepo{DB: db}

	authService := service.NewAuthService(userRepo, jwtSecret)
	taskService := service.NewTaskService(taskRepo, userTaskRepo)
	progressService := service.NewProgressService(db)

	authHandler := handler.NewAuthHandler(authService, userRepo)
	phaseHandler := handler.NewPhaseHandler(phaseRepo, weekRepo, dayRepo)
	taskHandler := handler.NewTaskHandler(taskService, taskRepo, userTaskRepo)
	progressHandler := handler.NewProgressHandler(progressService)

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		protected := api.Group("")
		protected.Use(middleware.Auth(jwtSecret))
		{
			protected.GET("/auth/me", authHandler.Me)
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

	return r
}
