package controllers

import (
	"net/http"
	"ozinshe/db/initializers"
	"ozinshe/internal/models"
	"ozinshe/internal/validations"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateMovie(c *gin.Context) {
	var userInput struct {
		NameOfProject string               `json:"nameOfProject" binding:"required,min=2"`
		CategoryID    uint                 `json:"categoryID" binding:"required"`
		TypeID        uint                 `json:"typeID" binding:"required"`
		AgeCategories []models.AgeCategory `json:"ageCategories" binding:"required"`
		Screenshots   []models.Screenshot  `json:"screenshots" binding:"required"`
		Seasons       []models.Season      `json:"seasons" binding:"required"`
		Year          string               `json:"year" binding:"required"`
		Timing        string               `json:"timing" binding:"required"`
		Keywords      string               `json:"keywords" binding:"required"`
		Description   string               `json:"description" binding:"required"`
		Director      string               `json:"director" binding:"required"`
		Producer      string               `json:"producer" binding:"required"`
		Cover         string               `json:"cover" binding:"required"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !validations.IsExistValue("categories", "id", userInput.CategoryID) ||
		!validations.IsExistValue("types", "id", userInput.TypeID) {

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"CategoryId": "The category or types does not exist!",
		})
		return
	}

	for _, ageCategory := range userInput.AgeCategories {
		if !validations.IsExistValue("age_categories", "id", ageCategory.ID) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"AgeCategoryId": "The age category does not exist!",
			})
			return
		}
	}

	if validations.IsUniqueValue("movies", "name_of_project", userInput.NameOfProject) {
		c.JSON(http.StatusConflict, gin.H{
			"Name": "The name of movie is already exist!",
		})
		return
	}

	movie := models.Movie{
		NameOfProject: userInput.NameOfProject,
		CategoryID:    userInput.CategoryID,
		TypeID:        userInput.TypeID,
		AgeCategories: userInput.AgeCategories,
		Screenshots:   userInput.Screenshots,
		Seasons:       userInput.Seasons,
		Year:          userInput.Year,
		Timing:        userInput.Timing,
		Keywords:      userInput.Keywords,
		Description:   userInput.Description,
		Director:      userInput.Director,
		Producer:      userInput.Producer,
		Cover:         userInput.Cover,
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
	result := initializers.DB.Preload("Screenshots").
		Preload("AgeCategories").
		Preload("Seasons", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Videos")
		}).
		First(&movie, id)

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
		NameOfProject string               `json:"nameOfProject" binding:"required,min=2"`
		CategoryID    uint                 `json:"categoryID" binding:"required"`
		TypeID        uint                 `json:"typeID" binding:"required"`
		AgeCategories []models.AgeCategory `json:"ageCategories" binding:"required"`
		Screenshots   []models.Screenshot  `json:"screenshots" binding:"required"`
		Seasons       []models.Season      `json:"seasons" binding:"required"`
		Year          string               `json:"year" binding:"required"`
		Timing        string               `json:"timing" binding:"required"`
		Keywords      string               `json:"keywords" binding:"required"`
		Description   string               `json:"description" binding:"required"`
		Director      string               `json:"director" binding:"required"`
		Producer      string               `json:"producer" binding:"required"`
		Cover         string               `json:"cover" binding:"required"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !validations.IsExistValue("categories", "id", userInput.CategoryID) ||
		!validations.IsExistValue("types", "id", userInput.TypeID) {

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"CategoryId": "The category or types does not exist!",
		})
		return
	}

	for _, ageCategory := range userInput.AgeCategories {
		if !validations.IsExistValue("age_categories", "id", ageCategory.ID) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"AgeCategoryId": "The age category does not exist!",
			})
			return
		}
	}

	if validations.IsUniqueValue("movies", "name_of_project", userInput.NameOfProject) {
		c.JSON(http.StatusConflict, gin.H{
			"Name": "The name of movie is already exist!",
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
		CategoryID:    userInput.CategoryID,
		TypeID:        userInput.TypeID,
		AgeCategories: userInput.AgeCategories,
		Screenshots:   userInput.Screenshots,
		Seasons:       userInput.Seasons,
		Year:          userInput.Year,
		Timing:        userInput.Timing,
		Keywords:      userInput.Keywords,
		Description:   userInput.Description,
		Director:      userInput.Director,
		Producer:      userInput.Producer,
		Cover:         userInput.Cover,
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

	err := initializers.DB.Delete(&movie).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	var favorite []models.Favorite

	resultFavorite := initializers.DB.Where("movie_id = ?", id).Find(&favorite)
	if err := resultFavorite.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	initializers.DB.Delete(&favorite)

	c.JSON(http.StatusOK, gin.H{
		"message": "The movie has been deleted successfully",
	})
}

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

func AgeCreateCategory(c *gin.Context) {
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
