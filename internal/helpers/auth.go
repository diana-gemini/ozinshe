package helpers

import (
	"net/http"

	"github.com/diana-gemini/ozinshe/api/middleware"

	"github.com/gin-gonic/gin"
)

func GetAuthUser(c *gin.Context) *middleware.AuthUser {
	authUser, exists := c.Get("authUser")

	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to get the user")
		return nil
	}

	if user, ok := authUser.(middleware.AuthUser); ok {
		return &user
	}

	return nil
}
