package middleware

import (
	"fmt"
	"net/http"
	"os"
	"ozinshe/db/initializers"
	"ozinshe/internal/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthUser struct {
	ID    uint   `json:"ID"`
	Email string `json:"Email"`
	Role  uint   `json:"Role"`
}

func RequireAuth(c *gin.Context) {
	var tokenString string
	tokenArray := strings.Split(c.GetHeader("Authorization"), " ")

	if len(tokenArray) >= 2 {
		tokenString = tokenArray[1]
		if tokenString == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		var user models.User
		initializers.DB.Find(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		authUser := AuthUser{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.RoleID,
		}

		c.Set("authUser", authUser)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		authUser, exists := c.Get("authUser")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Failed to get the user",
			})
			return
		}

		user, ok := authUser.(AuthUser)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		if user.Role != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
			})
			return
		}

		c.Next()
	}
}
