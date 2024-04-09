package controllers

import (
	"net/http"
	"ozinshe/db/initializers"
	"ozinshe/internal/models"
	"ozinshe/internal/validations"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

func CreateMovie(c *gin.Context) {
	var userInput struct {
		NameOfProject string          `json:"nameOfProject" binding:"required,min=2"`
		Category      string          `json:"category" binding:"required"`
		TypeOfProject string          `json:"typeOfProject" binding:"required"`
		AgeCategory   string          `json:"ageCategory" binding:"required"`
		Year          string          `json:"year" binding:"required"`
		Timing        string          `json:"timing" binding:"required"`
		Keywords      string          `json:"keywords" binding:"required"`
		Description   string          `json:"description" binding:"required"`
		Director      string          `json:"director" binding:"required"`
		Producer      string          `json:"producer" binding:"required"`
		CountOfSeason []models.Season `json:"countOfSeason" binding:"required"`
		Cover         string          `json:"cover" binding:"required"`
		Screenshots   pq.StringArray  `json:"screenshots" binding:"required"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if validations.IsUniqueValue("movies", "name_of_project", userInput.NameOfProject) {
		c.JSON(http.StatusConflict, gin.H{
			"Name": "The name of movie is already exist!",
		})
		return
	}

	movie := models.Movie{
		NameOfProject: userInput.NameOfProject,
		Category:      userInput.Category,
		TypeOfProject: userInput.TypeOfProject,
		AgeCategory:   userInput.AgeCategory,
		Year:          userInput.Year,
		Timing:        userInput.Timing,
		Keywords:      userInput.Keywords,
		Description:   userInput.Description,
		Director:      userInput.Director,
		Producer:      userInput.Producer,
		CountOfSeason: datatypes.NewJSONSlice(userInput.CountOfSeason),
		Cover:         userInput.Cover,
		Screenshots:   userInput.Screenshots,
	}

	result := initializers.DB.Create(&movie)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot create movie",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movie": movie,
	})
}

func EditMovie(c *gin.Context) {
	id := c.Param("id")

	var movie models.Movie
	result := initializers.DB.First(&movie, id)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"movie": "Record not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movie": movie,
	})
}

func UpdateMovie(c *gin.Context) {
	id := c.Param("id")

	var userInput struct {
		NameOfProject string          `json:"nameOfProject" binding:"required,min=2"`
		Category      string          `json:"category" binding:"required"`
		TypeOfProject string          `json:"typeOfProject" binding:"required"`
		AgeCategory   string          `json:"ageCategory" binding:"required"`
		Year          string          `json:"year" binding:"required"`
		Timing        string          `json:"timing" binding:"required"`
		Keywords      string          `json:"keywords" binding:"required"`
		Description   string          `json:"description" binding:"required"`
		Director      string          `json:"director" binding:"required"`
		Producer      string          `json:"producer" binding:"required"`
		CountOfSeason []models.Season `json:"countOfSeason" binding:"required"`
		Cover         string          `json:"cover" binding:"required"`
		Screenshots   pq.StringArray  `json:"screenshots" binding:"required"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var movie models.Movie
	result := initializers.DB.First(&movie, id)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	updateMovie := models.Movie{
		NameOfProject: userInput.NameOfProject,
		Category:      userInput.Category,
		TypeOfProject: userInput.TypeOfProject,
		AgeCategory:   userInput.AgeCategory,
		Year:          userInput.Year,
		Timing:        userInput.Timing,
		Keywords:      userInput.Keywords,
		Description:   userInput.Description,
		Director:      userInput.Director,
		Producer:      userInput.Producer,
		CountOfSeason: userInput.CountOfSeason,
		Cover:         userInput.Cover,
		Screenshots:   userInput.Screenshots,
	}

	result = initializers.DB.Model(&movie).Updates(&updateMovie)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movie": updateMovie,
	})
}

func DeleteMovie(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movie

	result := initializers.DB.First(&movie, id)
	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	initializers.DB.Delete(&movie)

	c.JSON(http.StatusOK, gin.H{
		"message": "The movie has been deleted successfully",
	})
}