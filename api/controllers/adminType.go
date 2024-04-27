package controllers

import (
	"net/http"
	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/models"
	"github.com/diana-gemini/ozinshe/internal/validations"

	"github.com/gin-gonic/gin"
)

func CreateTypeOfProject(c *gin.Context) {
	var userInput struct {
		TypeName string `json:"TypeName" binding:"required,min=2"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if validations.IsUniqueValue("types", "type_name", userInput.TypeName) {
		c.JSON(http.StatusConflict, gin.H{
			"validations": map[string]interface{}{
				"Name": "The type name is already exist!",
			},
		})
		return
	}

	typeOfProject := models.Type{
		TypeName: userInput.TypeName,
	}

	result := initializers.DB.Create(&typeOfProject)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot create category",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"typeOfProject": typeOfProject,
	})
}

func UpdateTypeOfProject(c *gin.Context) {
	id := c.Param("id")

	var userInput struct {
		TypeName string `json:"TypeName" binding:"required,min=2"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if validations.IsUniqueValue("types", "type_name", userInput.TypeName) {
		c.JSON(http.StatusConflict, gin.H{
			"validations": map[string]interface{}{
				"Name": "The type name is already exist!",
			},
		})
		return
	}

	var typeOfProject models.Type
	result := initializers.DB.First(&typeOfProject, id)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	updateTypeOfProject := models.Type{
		TypeName: userInput.TypeName,
	}

	result = initializers.DB.Model(&typeOfProject).Updates(&updateTypeOfProject)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot update type",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"typeName": updateTypeOfProject.TypeName,
	})
}

func DeleteTypeOfProject(c *gin.Context) {
	id := c.Param("id")
	var typeOfProject models.Type

	result := initializers.DB.First(&typeOfProject, id)
	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	err := initializers.DB.Delete(&typeOfProject).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "The type has been deleted successfully",
	})
}
