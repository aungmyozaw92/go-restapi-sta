package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aungmyozaw92/go-restapi-sta/cmd"
	"github.com/aungmyozaw92/go-restapi-sta/config"
	"github.com/aungmyozaw92/go-restapi-sta/models"
	"github.com/aungmyozaw92/go-restapi-sta/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = defaultPort
	}

	cmd.Execute()

	// Connect to Database
	db := config.GetDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	models.MigrateTable()

	// Initialize Gin router.
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE", "OPTIONS")
	corsConfig.AddAllowHeaders("token", "Origin", "Content-Type", "Authorization")
	corsConfig.AddExposeHeaders("Content-Length")
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))
	// Router
	routes.SetupRoutes(r)
	
	r.NoRoute(customNotFoundHandler)
	r.Run(":" + port)

	log.Println("Server started successfully")
}

func customNotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "route not found",
		})
}