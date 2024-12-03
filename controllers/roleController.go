package controllers

import (
	"net/http"
	"strconv"

	"github.com/aungmyozaw92/go-restapi-sta/models"
	"github.com/aungmyozaw92/go-restapi-sta/utils"
	"github.com/gin-gonic/gin"
)

func CreateRole(context *gin.Context) {
	var input models.NewRole

	if err := context.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
		return
	}

	responseData, err := models.CreateRole(context, &input)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "failed create role", err)
		return
	}
	// create role response
	utils.SuccessResponse(context, http.StatusOK, "role created", responseData)

}

func UpdateRole(context *gin.Context) {
	// Parse role ID from URL parameter
	idParam := context.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	var input models.NewRole
	if err := context.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "Invalid input", err)
		return
	}

	updatedRole, err := models.UpdateRole(context, id, &input)
	if err != nil {
		utils.ErrorResponse(context, http.StatusInternalServerError, "Failed to update role", err)
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "Role updated successfully", updatedRole)
}

func DeleteRole(context *gin.Context) {
	// Parse role ID from URL parameter
	idParam := context.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	deletedRole, err := models.DeleteRole(context, id)
	if err != nil {
		utils.ErrorResponse(context, http.StatusInternalServerError, "Failed to delete role", err)
		return
	}

	// Success response
	utils.SuccessResponse(context, http.StatusOK, "Role deleted successfully", deletedRole)
}

func GetRole(context *gin.Context) {
	// Parse role ID from URL parameter
	idParam := context.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	deletedRole, err := models.GetRole(context, id)
	if err != nil {
		utils.ErrorResponse(context, http.StatusInternalServerError, "Failed to get role", err)
		return
	}

	// Success response
	utils.SuccessResponse(context, http.StatusOK, "Success", deletedRole)
}


