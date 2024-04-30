package controllers

import (
	"net/http"
	"strconv"

	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/models"
	"github.com/diana-gemini/ozinshe/internal/validations"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NewMovie struct {
	NameOfProject string   `form:"nameOfProject" binding:"required" example:"Name of Film"`
	Categories    []string `form:"categories[]" binding:"required"`
	TypeID        string   `form:"typeID" binding:"required" example:"Film"`
	AgeCategoryID string   `form:"ageCategoryID" binding:"required" example:"13-17"`
	Year          string   `form:"year" binding:"required" example:"2001"`
	Timing        string   `form:"timing" binding:"required" example:"89"`
	Keywords      string   `form:"keywords" binding:"required" example:"Film, Horor, Anime"`
	Description   string   `form:"description" binding:"required" example:"Konec XX veka. Neskolko let nazad, eshyo devochkoj..."`
	Director      string   `form:"director" binding:"required" example:"Tomokadzu Tokoro, Hideki Tonokacu"`
	Producer      string   `form:"producer" binding:"required" example:"Satosi Fudzii, Yosiyuki Fudetani"`
}

// CreateMovie godoc
// @Summary CreateMovie
// @Security ApiKeyAuth
// @Tags admin-movie-controller
// @ID create-movie
// @Accept  json
// @Produce  json
// @Param id path integer true "movieID"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /movie/create [post]
func CreateMovie(c *gin.Context) {
	var newMovie NewMovie
	// nameOfProject := c.PostForm("nameOfProject")
	// categoriesArray := c.PostFormArray("categories")
	// typeID := c.PostForm("typeID")
	// ageCategoryID := c.PostForm("ageCategoryID")
	// year := c.PostForm("year")
	// timing := c.PostForm("timing")
	// keywords := c.PostForm("keywords")
	// description := c.PostForm("description")
	// director := c.PostForm("director")
	// producer := c.PostForm("producer")

	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	form := c.Request.MultipartForm

	screenshotsArray := form.File["screenshots"]
	cover := form.File["cover"]

	if validations.IsUniqueValue("movies", "name_of_project", newMovie.NameOfProject) {
		c.JSON(http.StatusConflict, gin.H{
			"Name": "The name of movie is already exist!",
		})
		return
	}

	var categories []models.Category

	for _, category := range newMovie.Categories {
		if !validations.IsExistValue("categories", "id", category) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"CategoryID": "The category does not exist!",
			})
			return
		}

		id, err := strconv.Atoi(category)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		var tempCategory models.Category
		if err = initializers.DB.First(&tempCategory, id).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		categories = append(categories, tempCategory)
	}

	if !validations.IsExistValue("types", "id", newMovie.TypeID) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"TypeID": "The type does not exist!",
		})
		return
	}

	typeIDInt, err := strconv.Atoi(newMovie.TypeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if !validations.IsExistValue("age_categories", "id", newMovie.AgeCategoryID) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"AgeCategoryID": "The age category does not exist!",
		})
		return
	}

	ageCategoryIDInt, err := strconv.Atoi(newMovie.AgeCategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	screenshotsURL, err := ImageUpload(c, screenshotsArray)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var screenshots []models.Screenshot

	for _, screenshotURL := range screenshotsURL {
		screenshot := models.Screenshot{
			Link: screenshotURL,
		}
		screenshots = append(screenshots, screenshot)
	}

	coverURL, err := ImageUpload(c, cover)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if len(coverURL) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cover is not found",
		})
		return
	}

	movie := models.Movie{
		NameOfProject: newMovie.NameOfProject,
		Categories:    categories,
		TypeID:        uint(typeIDInt),
		AgeCategoryID: uint(ageCategoryIDInt),
		Screenshots:   screenshots,
		Year:          newMovie.Year,
		Timing:        newMovie.Timing,
		Keywords:      newMovie.Keywords,
		Description:   newMovie.Description,
		Director:      newMovie.Director,
		Producer:      newMovie.Producer,
		Cover:         coverURL[0],
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

	nameOfProject := c.PostForm("nameOfProject")
	categoriesArray := c.PostFormArray("categories")
	typeID := c.PostForm("typeID")
	ageCategoryID := c.PostForm("ageCategoryID")
	//screenshotsArray := c.PostFormArray("screenshots")
	year := c.PostForm("year")
	timing := c.PostForm("timing")
	keywords := c.PostForm("keywords")
	description := c.PostForm("description")
	director := c.PostForm("director")
	producer := c.PostForm("producer")
	cover := c.PostForm("cover")

	if validations.IsUniqueValue("movies", "name_of_project", nameOfProject) {
		c.JSON(http.StatusConflict, gin.H{
			"Name": "The name of movie is already exist!",
		})
		return
	}

	var categories []models.Category

	for _, category := range categoriesArray {
		if !validations.IsExistValue("categories", "id", category) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"CategoryID": "The category does not exist!",
			})
			return
		}

		id, err := strconv.Atoi(category)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		var tempCategory models.Category
		if err = initializers.DB.First(&tempCategory, id).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		categories = append(categories, tempCategory)
	}

	if !validations.IsExistValue("types", "id", typeID) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"TypeID": "The type does not exist!",
		})
		return
	}

	typeIDInt, err := strconv.Atoi(typeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if !validations.IsExistValue("age_categories", "id", ageCategoryID) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"AgeCategoryID": "The age category does not exist!",
		})
		return
	}

	ageCategoryIDInt, err := strconv.Atoi(ageCategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
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
		NameOfProject: nameOfProject,
		Categories:    categories,
		TypeID:        uint(typeIDInt),
		AgeCategoryID: uint(ageCategoryIDInt),
		//Screenshots:   screenshots,
		Year:        year,
		Timing:      timing,
		Keywords:    keywords,
		Description: description,
		Director:    director,
		Producer:    producer,
		Cover:       cover,
	}

	if err := initializers.DB.Model(&movie).Association("Categories").Replace(updateMovie.Categories); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
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

	initializers.DB.Unscoped().Delete(&favorite)

	c.JSON(http.StatusOK, gin.H{
		"message": "The movie has been deleted successfully",
	})
}
