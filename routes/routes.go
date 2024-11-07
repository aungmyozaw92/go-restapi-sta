package routes

import (
	"net/http"

	"github.com/aungmyozaw92/go-restapi-sta/controllers"
	"github.com/aungmyozaw92/go-restapi-sta/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "welcome home")
	})

	authRouter := r.Group("/api/v1")
	authRouter.POST("/login", controllers.Login)
	
	protectedRouter := r.Group("/api/v1")
	protectedRouter.Use(middlewares.AuthMiddleware())

	protectedRouter.GET("/profile", controllers.Profile)
	// protectedRouter.POST("/logout", controllers.Logout)
}