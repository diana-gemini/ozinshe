package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/models"
	"github.com/diana-gemini/ozinshe/internal/validations"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NewMovie struct {
	NameOfProject string   `form:"nameOfProject" binding:"required" example:"Hellsing"`                      // field 1
	CategoriesID  []string `form:"categoriesID" binding:"required" example:"1"`                              // field 2
	TypeID        string   `form:"typeID" binding:"required" example:"1"`                                    // field 3
	AgeCategoryID string   `form:"ageCategoryID" binding:"required" example:"13-17"`                         // field 4
	Year          string   `form:"year" binding:"required" example:"2001"`                                   // field 5
	Timing        string   `form:"timing" binding:"required" example:"89"`                                   // field 6
	Keywords      string   `form:"keywords" binding:"required" example:"Film, Horor, Anime"`                 // field 7
	Description   string   `form:"description" binding:"required" example:"Konec XX veka. Neskolko let ..."` // field 8
	Director      string   `form:"director" binding:"required" example:"Tomokadzu Tokoro, Hideki Tonokacu"`  // field 9
	Producer      string   `form:"producer" binding:"required" example:"Satosi Fudzii, Yosiyuki Fudetani"`   // field 10
}

// CreateMovie godoc
// @Summary CreateMovie
// @Security ApiKeyAuth
// @Tags admin-movie-controller
// @ID create-movie
// @Accept multipart/form-data
// @Produce json
// @Param newMovie formData NewMovie true "newMovie"
// @Param screenshots formData []file true "screenshots" collectionFormat(multi) "Image files to upload"
// @Param seasons formData NewSeason true "seasons"
// @Param cover formData file true "cover"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /admin/movie/create [post]
func CreateMovie(c *gin.Context) {
	var newMovie NewMovie

	if err := c.ShouldBind(&newMovie); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if validations.IsUniqueValue("movies", "name_of_project", newMovie.NameOfProject) {
		NewErrorResponse(c, http.StatusConflict, "movie is already exist")
		return
	}

	categoriesArray := strings.Split(newMovie.CategoriesID[0], ",")

	var categories []models.Category
	for _, categoryID := range categoriesArray {
		if !validations.IsExistValue("categories", "id", categoryID) {
			NewErrorResponse(c, http.StatusBadRequest, "category does not exist")
			return
		}

		id, err := strconv.Atoi(categoryID)
		if err != nil {
			NewErrorResponse(c, http.StatusBadRequest, "cannot convert categogy ID to int")
			return
		}

		var tempCategory models.Category
		if err = initializers.DB.First(&tempCategory, id).Error; err != nil {
			NewErrorResponse(c, http.StatusBadRequest, "cannot find category")
			return
		}

		categories = append(categories, tempCategory)
	}

	if !validations.IsExistValue("types", "id", newMovie.TypeID) {
		NewErrorResponse(c, http.StatusBadRequest, "type does not exist")
		return
	}

	typeIDInt, err := strconv.Atoi(newMovie.TypeID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "cannot convert type ID to int")
		return
	}

	if !validations.IsExistValue("age_categories", "id", newMovie.AgeCategoryID) {
		NewErrorResponse(c, http.StatusBadRequest, "age category does not exist")
		return
	}

	ageCategoryIDInt, err := strconv.Atoi(newMovie.AgeCategoryID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "cannot convert age category ID to int")
		return
	}

	err = c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "failed to parse form")
		return
	}

	form := c.Request.MultipartForm

	screenshotsArray := form.File["screenshots"]

	screenshotsURL, err := ImageUpload(c, screenshotsArray)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "cannot file upload")
		return
	}

	var screenshots []models.Screenshot

	for _, screenshotURL := range screenshotsURL {
		screenshot := models.Screenshot{
			Link: screenshotURL,
		}
		screenshots = append(screenshots, screenshot)
	}

	cover := form.File["cover"]
	coverURL, err := ImageUpload(c, cover)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if len(coverURL) < 1 {
		NewErrorResponse(c, http.StatusBadRequest, "cover is not found")
		return
	}

	var newSeason NewSeason
	if err := c.ShouldBind(&newSeason); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	videosArray := strings.Split(newSeason.Videos[0], ",")

	var videos []models.Video

	for _, video := range videosArray {
		tempVideo := models.Video{
			Link: video,
		}
		videos = append(videos, tempVideo)
	}

	season := models.Season{
		Videos: videos,
	}
	var seasons []models.Season
	seasons = append(seasons, season)

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
		Seasons:       seasons,
	}

	result := initializers.DB.Create(&movie)

	if result.Error != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "cannot create movie")
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
		"message": "movie delete successfully",
	})
}
