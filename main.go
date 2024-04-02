package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var db *sql.DB

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost/mydb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()
	router.GET("/login", logIn)
	router.Run("localhost:8081")
}

func logIn(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}
