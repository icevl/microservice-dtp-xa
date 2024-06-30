package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const PORT = 8001

type NewUserPayload struct {
	Email string `json:"email" binding:"required"`
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	router := setUpRouter()
	fmt.Printf("Server running on port %d\n", PORT)

	err := router.Run(":" + fmt.Sprintf("%d", PORT))
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}

}

func setUpRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/create", newUser)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Route not found"})
	})

	return router
}

func newUser(context *gin.Context) {
	var payload NewUserPayload

	if err := context.BindJSON(&payload); err != nil {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "Invalid payload"})
	}

	id := uuid.New().String()

	context.JSON(http.StatusOK, gin.H{
		"message": payload.Email,
		"uuid":    id,
	})
}
