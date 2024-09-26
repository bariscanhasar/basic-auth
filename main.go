package main

import (
	"awesomeProject/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {

	r := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r.POST("/login", loginHandler)

	authorized := r.Group("/")
	authorized.Use(middleware.JWTAuth())
	{
		authorized.GET("/protected", protectedHandler)
	}

	err = r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

func loginHandler(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if user.Username == "admin" && user.Password == "password" {
		token, err := middleware.GenerateToken(user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

func protectedHandler(c *gin.Context) {
	userID, _ := c.Get("user_id")
	c.JSON(http.StatusOK, gin.H{"message": "This is a protected route", "user_id": userID})
}
