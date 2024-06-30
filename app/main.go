package main

import (
	"app/services"
	"app/types"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const PORT = 8001

func init() {
	loadEnvironmentVariables()
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	router := setupRouter()

	fmt.Printf("Server running on port %d\n", PORT)
	err := router.Run(":" + fmt.Sprintf("%d", PORT))
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}

}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/create", newUser)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Route not found"})
	})

	return router
}

func newUser(context *gin.Context) {
	var payload types.NewUserPayload

	if err := context.BindJSON(&payload); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Invalid payload"})
	}

	xid := uuid.New().String()

	var chain []services.Service
	chain = append(chain, &services.User{UUID: xid, Email: payload.Email})
	chain = append(chain, &services.Cash{UUID: xid})

	if response := preparingStage(chain); !response.Success {
		context.JSON(response.Status, gin.H{"success": false, "message": response.Message})
		return
	}

	if response := commitStage(chain); !response.Success {
		context.JSON(response.Status, gin.H{"success": false, "message": response.Message})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"success": true,
		"uuid":    xid,
	})

}

func loadEnvironmentVariables() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
}

func preparingStage(chain []services.Service) types.Result {
	resultCh := make(chan types.Result, 2)

	var wg sync.WaitGroup
	var response = types.Result{
		Status:  http.StatusOK,
		Success: true,
	}

	wg.Add(1)

	for _, service := range chain {
		go service.Create(resultCh)
	}

	go func() {
		defer wg.Done()

		for i := 0; i < len(chain); i++ {
			result := <-resultCh

			if !result.Success {
				// TODO: Sentry result.Error
				// Also can`t break here due queue of services to rollback
				response = types.Result{Message: result.Message, Status: result.Status, Success: false}
			}
		}
	}()

	wg.Wait()

	if !response.Success {
		rollback(chain)
	}

	return response
}

func commitStage(chain []services.Service) types.Result {
	resultCh := make(chan types.Result, 2)

	var wg sync.WaitGroup
	var response = types.Result{
		Status:  http.StatusOK,
		Success: true,
	}

	wg.Add(1)

	for _, service := range chain {
		go service.Commit(resultCh)
	}

	go func() {
		defer wg.Done()

		for i := 0; i < len(chain); i++ {
			result := <-resultCh

			if !result.Success {

				// TODO: Sentry result.Error

				response = types.Result{Message: result.Message, Status: result.Status, Success: false}
				break
			}
		}
	}()

	wg.Wait()

	return response
}

func rollback(chain []services.Service) {
	for _, service := range chain {
		go service.Rollback()
	}
}
