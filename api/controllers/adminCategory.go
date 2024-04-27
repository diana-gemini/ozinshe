package controllers

import (
	"net/http"
	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/models"
	"github.com/diana-gemini/ozinshe/internal/validations"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	var userInput struct {
		CategoryName string `json:"categoryName" binding:"required,min=2"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	if validations.IsUniqueValue("categories", "category_name", userInput.CategoryName) {
		c.JSON(http.StatusConflict, gin.H{
			"validations": map[string]interface{}{
				"Name": "The category name is already exist!",
			},
		})
		return
	}

	category := models.Category{
		CategoryName: userInput.CategoryName,
	}

	result := initializers.DB.Create(&category)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot create category",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": category,
	})
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var userInput struct {
		CategoryName string `json:"categoryName" binding:"required,min=2"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if validations.IsUniqueValue("categories", "category_name", userInput.CategoryName) {
		c.JSON(http.StatusConflict, gin.H{
			"validations": map[string]interface{}{
				"Name": "The category name is already exist!",
			},
		})
		return
	}

	var category models.Category
	result := initializers.DB.First(&category, id)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	updateCategory := models.Category{
		CategoryName: userInput.CategoryName,
	}

	result = initializers.DB.Model(&category).Updates(&updateCategory)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot update category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"UpdateCategoryName": updateCategory.CategoryName,
	})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	result := initializers.DB.First(&category, id)
	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	err := initializers.DB.Delete(&category).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "The category has been deleted successfully",
	})
}
