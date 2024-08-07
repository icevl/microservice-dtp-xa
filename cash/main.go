package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const PORT = 8003

func init() {
	loadEnvironmentVariables()
	gin.SetMode(gin.ReleaseMode)
}

func main() {

	DBConnect()

	router := setupRouter()
	fmt.Printf("Server running on port %d\n", PORT)

	err := router.Run(":" + fmt.Sprintf("%d", PORT))
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/create", create)
	router.POST("/commit", commit)
	router.POST("/rollback", rollback)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Route not found"})
	})

	return router
}

func loadEnvironmentVariables() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
}

func create(context *gin.Context) {
	var payload CreatePayload

	if err := context.BindJSON(&payload); err != nil {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "Invalid payload"})
	}

	err := StartTransaction(payload.UUID)
	if err != nil {
		failure(context, payload.UUID, ErrorResponse{Message: err.Error()})
		return
	}

	_, err = Database.Exec(fmt.Sprintf("INSERT INTO cash (uuid) VALUES ('%s')", payload.UUID))
	if err != nil {
		failure(context, payload.UUID, ErrorResponse{Message: err.Error()})
		return
	}

	err = EndTransaction(payload.UUID)
	if err != nil {
		fmt.Println(err)
		failure(context, payload.UUID, ErrorResponse{Message: err.Error()})
		return
	}

	err = PrepareTransaction(payload.UUID)
	if err != nil {
		fmt.Println(err)
		failure(context, payload.UUID, ErrorResponse{Message: err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func commit(context *gin.Context) {
	var payload OperationPayload

	if err := context.BindJSON(&payload); err != nil {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "Invalid payload"})
	}

	if err := CommitTransaction(payload.UUID); err != nil {
		failure(context, payload.UUID, ErrorResponse{Message: err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func rollback(context *gin.Context) {
	var payload OperationPayload

	if err := context.BindJSON(&payload); err != nil {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "Invalid payload"})
	}

	if err := RollbackTransaction(payload.UUID); err != nil {
		failure(context, payload.UUID, ErrorResponse{Message: err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func failure(context *gin.Context, uuid string, err ErrorResponse) {
	EndTransaction(uuid)
	RollbackTransaction(uuid)

	if err.Code == 0 {
		err.Code = http.StatusInternalServerError
	}

	context.JSON(err.Code, gin.H{
		"success": false,
		"message": err.Message,
	})
}
