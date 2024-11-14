package controllers

import (
	"net/http"

	"github.com/aungmyozaw92/go-restapi-sta/models"
	"github.com/aungmyozaw92/go-restapi-sta/utils"
	"github.com/gin-gonic/gin"
)


type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(context *gin.Context) {
	var input LoginInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	responseData, err := models.Login(context, input.Username, input.Password)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	// Successful login response
	context.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Login successful",
		"data": responseData,
	})
}

func Profile(context *gin.Context) {
   	userId, ok := utils.GetUserIdFromContext(context.Request.Context())

	if !ok || userId == 0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "user id is required",
		})
		return
	}

	responseData, err := models.GetUser(context,userId)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Successful login response
	context.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "success",
		"data": responseData,
	})
}