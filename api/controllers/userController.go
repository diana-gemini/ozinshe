package controllers

import (
	"net/http"
	"os"
	"time"

	"ozinshe/db/initializers"
	format_errors "ozinshe/internal/format-errors"
	"ozinshe/internal/models"
	"ozinshe/internal/validations"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var userInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=4"`
	}

	userInput.Email = c.PostForm("email")
	userInput.Password = c.PostForm("password")
	errs := models.ErrText{}

	if validations.IsUniqueValue("users", "email", userInput.Email) {
		errs.Email = "Email already exists"
		c.HTML(http.StatusBadRequest, "signup.html", errs)
		return
	}

	if !CheckPassword(c) {
		errs.Pass2 = "Passwords should be the same"
		c.HTML(http.StatusBadRequest, "signup.html", errs)
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
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		format_errors.InternalServerError(c)
		return
	}
	c.Redirect(http.StatusSeeOther, "/login")
}

func Login(c *gin.Context) {
	var userInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	userInput.Email = c.PostForm("email")
	userInput.Password = c.PostForm("password")

	var user models.User
	errs := models.ErrText{}
	initializers.DB.First(&user, "email = ?", userInput.Email)

	if user.ID == 0 {
		errs.Email = "email not found"
		c.HTML(http.StatusBadRequest, "login.html", errs)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	if err != nil {
		errs.Email = "Invalid email or password"
		c.HTML(http.StatusBadRequest, "login.html", errs)
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

	auth := models.WebPage{}
	auth.IsLoggedin = true

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
	c.HTML(http.StatusOK, "index.html", auth)
}

func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", 0, "", "", false, true)
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func CheckPassword(c *gin.Context) bool {
	password := c.PostForm("password")
	if password == "" {
		return false
	}
	if len(password) < 4 {
		return false
	}
	if len(password) > 50 {
		return false
	}
	if PasswordRepeat(c, password) {
		return true
	}
	return false
}

func PasswordRepeat(c *gin.Context, firstPass string) bool {
	repeatPassword := c.PostForm("password2")
	if repeatPassword == "" {
		return false
	}
	if firstPass != repeatPassword {
		return false
	}
	return true
}
