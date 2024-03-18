package router

import (
	"social-media-app/controllers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartApp(DB *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("DB", DB)
		c.Next()
	})

	router.GET("/health-check", controllers.ServerCheck)

	return router
}
