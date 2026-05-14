package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xyd/web3-learning-tracker/internal/middleware"
)

func Setup(jwtSecret string) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	api := r.Group("/api/v1")
	api.Use(middleware.Auth(jwtSecret))
	{
		_ = api
	}

	return r
}
