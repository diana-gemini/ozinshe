package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {
	var err error

	DB := "user=postgres_user password=postgres_password dbname=postgres_db port=5432 sslmode=disable"

	db, err = sql.Open("postgres", DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title": "Home Page",
			},
		)
	})
	router.POST("/signup", signup)
	router.Run("localhost:8080")
}

func signup(c *gin.Context) {
	// c.IndentedJSON(http.StatusOK, albums)
	var user struct {
		Email    string
		Password string
	}

	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")

	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	result, err := db.Exec("insert into Users (email, password) values ($1, $2)",
		user.Email, string(hash))
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
	// user := models.User{Email: body.Email, Password: string(hash)}
}
