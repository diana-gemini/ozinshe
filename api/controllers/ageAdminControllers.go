package controllers

import (
	"net/http"
	"ozinshe/db/initializers"
	"ozinshe/internal/models"
	"ozinshe/internal/validations"

	"github.com/gin-gonic/gin"
)

func CreateAgeCategory(c *gin.Context) {
	var userInput struct {
		AgeCategoryName string `json:"ageCategoryName" binding:"required,min=2"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if validations.IsUniqueValue("age_categories", "age_category_name", userInput.AgeCategoryName) {
		c.JSON(http.StatusConflict, gin.H{
			"validations": map[string]interface{}{
				"Name": "The category name is already exist!",
			},
		})
		return
	}

	ageCategory := models.AgeCategory{
		AgeCategoryName: userInput.AgeCategoryName,
	}

	result := initializers.DB.Create(&ageCategory)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot create category",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ageCategory": ageCategory,
	})
}

func UpdateAgeCategory(c *gin.Context) {
	id := c.Param("id")

	var userInput struct {
		AgeCategoryName string `json:"ageCategoryName" binding:"required,min=2"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if validations.IsUniqueValue("age_categories", "age_category_name", userInput.AgeCategoryName) {
		c.JSON(http.StatusConflict, gin.H{
			"validations": map[string]interface{}{
				"Name": "The category name is already exist!",
			},
		})
		return
	}

	var ageCategory models.AgeCategory
	result := initializers.DB.First(&ageCategory, id)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	updateAgeCategory := models.AgeCategory{
		AgeCategoryName: userInput.AgeCategoryName,
	}

	result = initializers.DB.Model(&ageCategory).Updates(&updateAgeCategory)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot update age category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ageCategory": updateAgeCategory.AgeCategoryName,
	})
}

func DeleteAgeCategory(c *gin.Context) {
	id := c.Param("id")
	var ageCategory models.AgeCategory

	result := initializers.DB.First(&ageCategory, id)
	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	err := initializers.DB.Delete(&ageCategory).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "The age category has been deleted successfully",
	})
}
