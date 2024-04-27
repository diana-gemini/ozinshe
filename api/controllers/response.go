package controllers

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Message string `json:"message"`
}

// type statusResponse struct {
// 	Status string `json:"status"`
// }

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, ErrorResponse{message})
}
