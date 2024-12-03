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
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
		return
	}

	responseData, err := models.Login(context, input.Username, input.Password)
	if err != nil {
		utils.ErrorResponse(context, http.StatusUnauthorized, err.Error(), nil)
		return
	}
	// Successful login response
	
	utils.SuccessResponse(context, http.StatusOK, "Login successful", responseData)

}

func Profile(context *gin.Context) {
   	userId, ok := utils.GetUserIdFromContext(context.Request.Context())
	if !ok || userId == 0 {
		utils.ErrorResponse(context, http.StatusBadRequest, "user id is required", nil)
		return
	}

	responseData, err := models.GetUser(context,userId)

	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Successful profile response
	utils.SuccessResponse(context, http.StatusOK, "success", responseData)
}