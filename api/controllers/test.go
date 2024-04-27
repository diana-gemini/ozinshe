package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// type Video struct {
// 	Link string `form:"link"`
// }

type Video struct {
	Link string `form:"link"`
}

type Season struct {
	Videos []Video `form:"videos"`
}

func TestSeason(c *gin.Context) {
	var requestData Season

	// Parse form data
	if err := c.ShouldBind(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("requestData - %#v \n", requestData)
}
