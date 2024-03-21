package router

import (
	"social-media-app/controllers"
	"social-media-app/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartApp(DB *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("DB", DB)
		c.Next()
	})

	postRouter := router.Group("v1/post")
	{
		postRouter.POST("/", middleware.PostValidator(), controllers.CreatePost)
		postRouter.GET("/", controllers.GetPost)
	}

	router.GET("/seed-test", controllers.CreateSeed)

	router.GET("/health-check", controllers.ServerCheck)

	return router
}
