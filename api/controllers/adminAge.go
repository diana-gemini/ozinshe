package controllers

import (
	"net/http"

	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/models"
	"github.com/diana-gemini/ozinshe/internal/validations"

	"github.com/gin-gonic/gin"
)

type NewAgeCategory struct {
	AgeCategoryName string `json:"ageCategoryName" binding:"required,min=2" example:"0-13"`
}

// CreateAgeCategory godoc
// @Summary CreateAgeCategory
// @Security ApiKeyAuth
// @Tags admin-movie-age-category-controller
// @ID create-age-category
// @Accept  json
// @Produce  json
// @Param ageCategoryName body NewAgeCategory true "ageCategoryName"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /admin/agecategory/create [post]
func CreateAgeCategory(c *gin.Context) {
	var userInput NewAgeCategory

	if err := c.ShouldBindJSON(&userInput); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if validations.IsUniqueValue("age_categories", "age_category_name", userInput.AgeCategoryName) {
		NewErrorResponse(c, http.StatusConflict, "age category name is already exist")
		return
	}

	ageCategory := models.AgeCategory{
		AgeCategoryName: userInput.AgeCategoryName,
	}

	result := initializers.DB.Create(&ageCategory)

	if result.Error != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "cannot create age category")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ageCategory": ageCategory,
	})
}

// EditAgeCategory godoc
// @Summary EditAgeCategory
// @Security ApiKeyAuth
// @Tags admin-movie-age-category-controller
// @ID edit-age-category
// @Accept  json
// @Produce  json
// @Param id path integer true "id"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /admin/agecategory/{id}/edit [get]
func EditAgeCategory(c *gin.Context) {
	id := c.Param("id")

	var ageCategory models.AgeCategory
	result := initializers.DB.First(&ageCategory, id)

	if err := result.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "age category not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ageCategory": ageCategory,
	})
}

// UpdateAgeCategory godoc
// @Summary UpdateAgeCategory
// @Security ApiKeyAuth
// @Tags admin-movie-age-category-controller
// @ID update-age-category
// @Accept  json
// @Produce  json
// @Param id path integer true "id"
// @Param ageCategoryName body NewAgeCategory true "ageCategoryName"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /admin/agecategory/{id}/update [put]
func UpdateAgeCategory(c *gin.Context) {
	id := c.Param("id")

	var userInput NewAgeCategory

	if err := c.ShouldBindJSON(&userInput); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	var ageCategory models.AgeCategory
	result := initializers.DB.First(&ageCategory, id)

	if err := result.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "cannot find age category")
		return
	}

	if validations.IsUniqueValue("age_categories", "age_category_name", userInput.AgeCategoryName) {
		NewErrorResponse(c, http.StatusConflict, "age category name is already exist")
		return
	}

	updateAgeCategory := models.AgeCategory{
		AgeCategoryName: userInput.AgeCategoryName,
	}

	result = initializers.DB.Model(&ageCategory).Updates(&updateAgeCategory)

	if result.Error != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "cannot update age category")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ageCategory": updateAgeCategory.AgeCategoryName,
	})
}

// DeleteAgeCategory godoc
// @Summary DeleteAgeCategory
// @Security ApiKeyAuth
// @Tags admin-movie-age-category-controller
// @ID delete-age-category
// @Accept  json
// @Produce  json
// @Param id path integer true "id"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /admin/agecategory/{id}/delete [delete]
func DeleteAgeCategory(c *gin.Context) {
	id := c.Param("id")
	var ageCategory models.AgeCategory

	result := initializers.DB.First(&ageCategory, id)
	if err := result.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "cannot find age category")
		return
	}

	err := initializers.DB.Delete(&ageCategory).Error
	if err != nil {
		NewErrorResponse(c, http.StatusNotFound, "cannot delete age category")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "age category delete successfully",
	})
}
