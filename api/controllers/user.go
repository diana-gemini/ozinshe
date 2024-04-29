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

type SignUpUser struct {
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
// @Param signupRequest body SignUpUser true "signupRequest"
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /signup [post]
func Signup(c *gin.Context) {
	var userInput SignUpUser

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

type AuthUser struct {
	Email    string `json:"email" binding:"required,email" example:"user@mail.ru"`
	Password string `json:"password" binding:"required" example:"123456789"`
}

// Login godoc
// @Summary Login
// @Tags auth-controller
// @ID log-in
// @Accept  json
// @Produce  json
// @Param loginRequest body AuthUser true "loginRequest"
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /login [post]
func Login(c *gin.Context) {
	var userInput AuthUser

	if err := c.ShouldBindJSON(&userInput); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", userInput.Email)
	if user.ID == 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid email")
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "failed to create token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Token": tokenString,
	})
}

// Logout godoc
// @Summary Logout
// @Security ApiKeyAuth
// @Tags auth-controller
// @ID logout
// @Accept  json
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /logout [post]
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Logout": "Logout successful",
	})
}

// EditUserProfile godoc
// @Summary EditUserProfile
// @Security ApiKeyAuth
// @Tags user-controller
// @ID edit-user-profile
// @Accept  json
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /editprofile [get]
func EditUserProfile(c *gin.Context) {
	authUser := helpers.GetAuthUser(c)
	id := authUser.ID

	var user models.User
	result := initializers.DB.Select("username", "email", "mobile_phone", "birth_date").First(&user, id)
	if err := result.Error; err != nil {
		newErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"username":    user.Username,
		"email":       user.Email,
		"mobilePhone": user.MobilePhone,
		"birthDate":   user.BirthDate,
	})
}

type UserProfile struct {
	Username    string `json:"username" binding:"min=2" example:"Tilda"`
	MobilePhone string `json:"mobilePhone" example:"+7(705)1112233"`
	BirthDate   string `json:"birthDate" example:"01.01.2000"`
}

// UpdateUserProfile godoc
// @Summary UpdateUserProfile
// @Security ApiKeyAuth
// @Tags user-controller
// @ID update-user-profile
// @Accept  json
// @Produce  json
// @Param updateUserProfile body UserProfile true "updateUserProfile"
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /updateprofile [put]
func UpdateUserProfile(c *gin.Context) {
	authUser := helpers.GetAuthUser(c)
	id := authUser.ID

	var userInput UserProfile

	if err := c.ShouldBindJSON(&userInput); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	var user models.User
	result := initializers.DB.First(&user, id)

	if err := result.Error; err != nil {
		newErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	updateUser := models.User{
		Username:    userInput.Username,
		MobilePhone: userInput.MobilePhone,
		BirthDate:   userInput.BirthDate,
	}

	result = initializers.DB.Model(&user).Updates(&updateUser)

	if result.Error != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user data is not updated")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":    updateUser.Username,
		"mobilePhone": updateUser.MobilePhone,
		"birthDate":   updateUser.BirthDate,
	})
}

type UserPassword struct {
	Password       string `json:"password" binding:"required" example:"123456789"`
	RepeatPassword string `json:"passwordrepeat" binding:"required" example:"123456789"`
}

// ChangePassword godoc
// @Summary ChangePassword
// @Security ApiKeyAuth
// @Tags user-controller
// @ID change-password
// @Accept  json
// @Produce  json
// @Param changePassword body UserPassword true "changePassword"
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /changepassword [post]
func ChangePassword(c *gin.Context) {
	authUser := helpers.GetAuthUser(c)
	id := authUser.ID

	var userInput UserPassword

	if err := c.ShouldBindJSON(&userInput); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
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

	var user models.User
	result := initializers.DB.First(&user, id)

	if err := result.Error; err != nil {
		newErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	updateUserPassword := models.User{
		Password: string(hashPassword),
	}

	result = initializers.DB.Model(&user).Updates(&updateUserPassword)

	if result.Error != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user password is not changed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"password": "Password successfully changed",
	})
}

var resetTokens = make(map[string]string)

type RecoverUserPassword struct {
	Email string `json:"email" binding:"required,email" example:"user@mail.ru"`
}

// PasswordRecover godoc
// @Summary PasswordRecover
// @Tags password-controller
// @ID password-recover
// @Accept  json
// @Produce  json
// @Param email body RecoverUserPassword true "email"
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /passwordrecover [post]
func PasswordRecover(c *gin.Context) {
	var userInput RecoverUserPassword

	if err := c.ShouldBindJSON(&userInput); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", userInput.Email)

	if user.ID == 0 {
		newErrorResponse(c, http.StatusBadRequest, "user not found")
		return
	}

	token := generateToken()
	resetTokens[token] = userInput.Email

	err := sendResetEmail(userInput.Email, token)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to send reset email")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset email sent",
	})
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
	m.SetBody("text/html", fmt.Sprintf("Click the following link to reset your password: <a href=\"http://localhost:3000/resetpassword?token=%s\">Reset Password</a>", token))

	d := gomail.NewDialer("smtp.mail.ru", 587, "diana-test-project@mail.ru", "phwEkPnEPudpnmU0Pvvn")

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// ResetPassword godoc
// @Summary ResetPassword
// @Tags password-controller
// @ID password-reset
// @Accept  json
// @Produce  json
// @Param token query string true "token received in the URL"
// @Param resetPassword body UserPassword true "resetPassword"
// @Success 200 {integer} integer 1
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /resetpassword [post]
func ResetPassword(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		newErrorResponse(c, http.StatusBadRequest, "token not found in URL")
		return
	}

	var password UserPassword
	if err := c.ShouldBindJSON(&password); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	email, ok := resetTokens[token]
	if !ok {
		newErrorResponse(c, http.StatusBadRequest, "invalid reset token")
		return
	}

	if !validations.CheckPassword(password.Password, password.RepeatPassword) {
		newErrorResponse(c, http.StatusBadRequest, "passwords should be the same")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password.Password), bcrypt.DefaultCost)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to generate password hash")
		return
	}

	delete(resetTokens, token)

	var user models.User
	result := initializers.DB.Where("email = ?", email).First(&user)

	if err := result.Error; err != nil {
		newErrorResponse(c, http.StatusNotFound, "email not found")
		return
	}

	updateUserPassword := models.User{
		Password: string(hashedPassword),
	}

	result = initializers.DB.Model(&user).Updates(&updateUserPassword)

	if result.Error != nil {
		newErrorResponse(c, http.StatusInternalServerError, "password not update")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"password": "Password successfully changed",
	})
}
