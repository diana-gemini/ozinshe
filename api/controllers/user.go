package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/diana-gemini/ozinshe/db/initializers"
	"github.com/diana-gemini/ozinshe/internal/helpers"
	"github.com/diana-gemini/ozinshe/internal/models"
	"github.com/diana-gemini/ozinshe/internal/validations"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type AuthUser struct {
	Email          string `json:"email" binding:"required,email" example:"user@mail.ru"`
	Password       string `json:"password" binding:"required" example:"123456789"`
	RepeatPassword string `json:"passwordrepeat" binding:"required" example:"123456789"`
}

// Signup godoc
// @Summary SignUp
// @Tags auth-controller
// @ID create-account
// @Accept  json
// @Produce  json
// @Param signupRequest body AuthUser true "signupRequest"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /signup [post]
func Signup(c *gin.Context) {
	var userInput AuthUser

	if err := c.ShouldBindJSON(&userInput); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if validations.IsUniqueValue("users", "email", userInput.Email) {
		newErrorResponse(c, http.StatusBadRequest, "email is already exist")
		return
	}

	if !validations.CheckPassword(userInput.Password, userInput.RepeatPassword) {
		newErrorResponse(c, http.StatusBadRequest, "passwords should be the same")
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), 10)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user := models.User{
		Email:    userInput.Email,
		Password: string(hashPassword),
		RoleID:   2,
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"User": user,
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
	authUser := helpers.GetAuthUser(c)
	id := authUser.ID
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
	authUser := helpers.GetAuthUser(c)
	id := authUser.ID

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
	authUser := helpers.GetAuthUser(c)
	id := authUser.ID

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

var resetTokens = make(map[string]string)

func PasswordRecover(c *gin.Context) {
	email := c.PostForm("email")

	var user models.User
	initializers.DB.First(&user, "email = ?", email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	token := generateToken()
	resetTokens[token] = email

	err := sendResetEmail(email, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func sendResetEmail(email, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "diana-test-project@mail.ru")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Password Reset")
	m.SetBody("text/html", fmt.Sprintf("Click the following link to reset your password: <a href=\"http://localhost:3000/reset?token=%s\">Reset Password</a>", token))

	d := gomail.NewDialer("smtp.mail.ru", 587, "diana-test-project@mail.ru", "phwEkPnEPudpnmU0Pvvn")

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func ResetPassword(c *gin.Context) {
	token := c.PostForm("token")
	password := c.PostForm("password")
	passwordRepeat := c.PostForm("passwordrepeat")

	email, ok := resetTokens[token]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reset token"})
		return
	}

	if !validations.CheckPassword(password, passwordRepeat) {
		c.JSON(http.StatusBadRequest, gin.H{
			"Password": "Passwords should be the same",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate password hash"})
		return
	}

	delete(resetTokens, token)

	var user models.User
	result := initializers.DB.Where("email = ?", email).First(&user)

	if err := result.Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	updateUserPassword := models.User{
		Password: string(hashedPassword),
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
