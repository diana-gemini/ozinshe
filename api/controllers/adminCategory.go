package controllers

import (
	"net/http"

	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/models"
	"github.com/diana-gemini/ozinshe/internal/validations"

	"github.com/gin-gonic/gin"
)

type NewCategory struct {
	CategoryName string `json:"categoryName" binding:"required,min=2" example:"Anime"`
}

// CreateCategory godoc
// @Summary CreateCategory
// @Security ApiKeyAuth
// @Tags admin-movie-category-controller
// @ID create-category
// @Accept  json
// @Produce  json
// @Param categoryName body NewCategory true "categoryName"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /admin/category/create [post]
func CreateCategory(c *gin.Context) {
	var userInput NewCategory

	if err := c.ShouldBindJSON(&userInput); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if validations.IsUniqueValue("categories", "category_name", userInput.CategoryName) {
		NewErrorResponse(c, http.StatusConflict, "the category name is already exist")
		return
	}

	category := models.Category{
		CategoryName: userInput.CategoryName,
	}

	result := initializers.DB.Create(&category)

	if result.Error != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "cannot create category")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": category,
	})
}

// EditCategory godoc
// @Summary EditCategory
// @Security ApiKeyAuth
// @Tags admin-movie-category-controller
// @ID edit-category
// @Accept  json
// @Produce  json
// @Param id path integer true "id"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /admin/category/{id}/edit [get]
func EditCategory(c *gin.Context) {
	id := c.Param("id")

	var category models.Category
	result := initializers.DB.First(&category, id)

	if err := result.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "category not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": category,
	})
}

// UpdateCategory godoc
// @Summary UpdateCategory
// @Security ApiKeyAuth
// @Tags admin-movie-category-controller
// @ID update-category
// @Accept  json
// @Produce  json
// @Param id path integer true "id"
// @Param categoryName body NewCategory true "categoryName"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /admin/category/{id}/update [put]
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var userInput NewCategory

	if err := c.ShouldBindJSON(&userInput); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	var category models.Category
	result := initializers.DB.First(&category, id)

	if err := result.Error; err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "cannot find category")
		return
	}

	if validations.IsUniqueValue("categories", "category_name", userInput.CategoryName) {
		NewErrorResponse(c, http.StatusConflict, "the category name is already exist")
		return
	}

	updateCategory := models.Category{
		CategoryName: userInput.CategoryName,
	}

	result = initializers.DB.Model(&category).Updates(&updateCategory)

	if result.Error != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "cannot update category")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"UpdateCategoryName": updateCategory.CategoryName,
	})
}

// DeleteCategory godoc
// @Summary DeleteCategory
// @Security ApiKeyAuth
// @Tags admin-movie-category-controller
// @ID delete-category
// @Accept  json
// @Produce  json
// @Param id path integer true "id"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /admin/category/{id}/delete [delete]
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	result := initializers.DB.First(&category, id)
	if err := result.Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "cannot find category")
		return
	}

	err := initializers.DB.Delete(&category).Error
	if err != nil {
		NewErrorResponse(c, http.StatusNotFound, "cannot delete category")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "category delete successfully",
	})
}
