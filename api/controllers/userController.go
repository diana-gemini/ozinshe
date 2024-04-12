package controllers

import (
	"net/http"
	"os"
	"time"

	"ozinshe/api/middleware"
	"ozinshe/db/initializers"
	"ozinshe/internal/models"
	"ozinshe/internal/validations"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var userInput struct {
		Email          string `json:"email" binding:"required,email"`
		Password       string `json:"password" binding:"required,min=4"`
		RepeatPassword string `json:"passwordrepeat" binding:"required,min=4"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if validations.IsUniqueValue("users", "email", userInput.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"Email": "The email is already exist!",
		})
		return
	}

	if !validations.CheckPassword(userInput.Password, userInput.RepeatPassword) {
		c.JSON(http.StatusBadRequest, gin.H{
			"Password": "Passwords should be the same",
		})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{
		Email:    userInput.Email,
		Password: string(hashPassword),
		RoleID:   2,
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func Login(c *gin.Context) {
	var userInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if c.ShouldBindJSON(&userInput) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", userInput.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"successMessage": "Logout successful",
	})
}

func EditUserProfile(c *gin.Context) {
	id := middleware.GetAuthUserID(c)
	var user models.User
	result := initializers.DB.Select("username", "email", "mobile_phone", "birth_date").First(&user, id)
	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"user": "User not found",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"username":    user.Username,
		"email":       user.Email,
		"mobilePhone": user.MobilePhone,
		"birthDate":   user.BirthDate,
	})
}

func UpdateUserProfile(c *gin.Context) {
	id := middleware.GetAuthUserID(c)

	var userInput struct {
		Username    string `json:"username" binding:"min=2"`
		MobilePhone string `json:"mobilePhone"`
		BirthDate   string `json:"birthDate" `
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, id)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	updateUser := models.User{
		Username:    userInput.Username,
		MobilePhone: userInput.MobilePhone,
		BirthDate:   userInput.BirthDate,
	}

	result = initializers.DB.Model(&user).Updates(&updateUser)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":    updateUser.Username,
		"mobilePhone": updateUser.MobilePhone,
		"birthDate":   updateUser.BirthDate,
	})
}

func ChangePassword(c *gin.Context) {
	id := middleware.GetAuthUserID(c)

	var userInput struct {
		Password       string `json:"password" binding:"required,min=4"`
		RepeatPassword string `json:"passwordrepeat" binding:"required,min=4"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !validations.CheckPassword(userInput.Password, userInput.RepeatPassword) {
		c.JSON(http.StatusBadRequest, gin.H{
			"Password": "Passwords should be the same",
		})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, id)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	updateUserPassword := models.User{
		Password: string(hashPassword),
	}

	result = initializers.DB.Model(&user).Updates(&updateUserPassword)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"password": "Password successfully changed",
	})
}
