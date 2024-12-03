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

	// role routes
	roleRouter := protectedRouter.Group("/roles")
	{
		roleRouter.POST("", controllers.CreateRole)       // Create role
		roleRouter.GET("/:id", controllers.GetRole)       // Get Role by ID
		roleRouter.PUT("/:id", controllers.UpdateRole)    // Update Role by ID
		roleRouter.DELETE("/:id", controllers.DeleteRole) // Delete Role by ID
	// 	// roleRouter.GET("", controllers.ListRoles)         // List all Roles
	}

	// User routes
	// userRouter := protectedRouter.Group("/users")
	// {
	// 	// userRouter.POST("", controllers.CreateUser)       // Create user
	// 	// userRouter.GET("/:id", controllers.GetUser)       // Get user by ID
	// 	// userRouter.PUT("/:id", controllers.UpdateUser)    // Update user by ID
	// 	// userRouter.DELETE("/:id", controllers.DeleteUser) // Delete user by ID
	// 	// userRouter.GET("", controllers.ListUsers)         // List all users
	// }

}