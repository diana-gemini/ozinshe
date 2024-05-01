package controllers

import (
	"net/http"

	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/models"
	"github.com/diana-gemini/ozinshe/internal/validations"

	"github.com/gin-gonic/gin"
)

type NewType struct {
	TypeName string `json:"typeName" binding:"required,min=2" example:"Serial"`
}

// CreateTypeOfProject godoc
// @Summary CreateTypeOfProject
// @Security ApiKeyAuth
// @Tags admin-movie-type-controller
// @ID create-type
// @Accept  json
// @Produce  json
// @Param typeName body NewType true "typeName"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /admin/type/create [post]
func CreateTypeOfProject(c *gin.Context) {
	var userInput NewType

	if err := c.ShouldBindJSON(&userInput); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if validations.IsUniqueValue("types", "type_name", userInput.TypeName) {
		NewErrorResponse(c, http.StatusConflict, "the type name is already exist")
		return
	}

	typeOfProject := models.Type{
		TypeName: userInput.TypeName,
	}

	result := initializers.DB.Create(&typeOfProject)

	if result.Error != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "cannot create category")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"typeOfProject": typeOfProject,
	})
}

// EditTypeOfProject godoc
// @Summary EditTypeOfProject
// @Security ApiKeyAuth
// @Tags admin-movie-type-controller
// @ID edit-type
// @Accept  json
// @Produce  json
// @Param id path integer true "id"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /admin/type/{id}/edit [get]
func EditTypeOfProject(c *gin.Context) {
	id := c.Param("id")

	var typeOfProject models.Type
	result := initializers.DB.First(&typeOfProject, id)

	if err := result.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "type not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type": typeOfProject,
	})
}

// UpdateTypeOfProject godoc
// @Summary UpdateTypeOfProject
// @Security ApiKeyAuth
// @Tags admin-movie-type-controller
// @ID update-type
// @Accept  json
// @Produce  json
// @Param id path integer true "id"
// @Param typeName body NewType true "typeName"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /admin/type/{id}/update [put]
func UpdateTypeOfProject(c *gin.Context) {
	id := c.Param("id")

	var userInput NewType

	if err := c.ShouldBindJSON(&userInput); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if validations.IsUniqueValue("types", "type_name", userInput.TypeName) {
		NewErrorResponse(c, http.StatusConflict, "the type name is already exist")
		return
	}

	var typeOfProject models.Type
	result := initializers.DB.First(&typeOfProject, id)

	if err := result.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "cannot find type")
		return
	}

	updateTypeOfProject := models.Type{
		TypeName: userInput.TypeName,
	}

	result = initializers.DB.Model(&typeOfProject).Updates(&updateTypeOfProject)

	if result.Error != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "cannot update type")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"typeName": updateTypeOfProject.TypeName,
	})
}

// DeleteTypeOfProject godoc
// @Summary DeleteTypeOfProject
// @Security ApiKeyAuth
// @Tags admin-movie-type-controller
// @ID delete-type
// @Accept  json
// @Produce  json
// @Param id path integer true "id"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /admin/type/{id}/delete [delete]
func DeleteTypeOfProject(c *gin.Context) {
	id := c.Param("id")
	var typeOfProject models.Type

	result := initializers.DB.First(&typeOfProject, id)
	if err := result.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "cannot find type")
		return
	}

	err := initializers.DB.Delete(&typeOfProject).Error
	if err != nil {
		NewErrorResponse(c, http.StatusNotFound, "cannot find type")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "type delete successfully",
	})
}
