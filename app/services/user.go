package services

import (
	"app/types"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type User struct {
	UUID  string `json:"uuid"`
	Email string `json:"email"`
}

var userServicePath = loadUserServicePath()

func loadUserServicePath() string {
	godotenv.Load(".env")
	return os.Getenv("USER_SERVICE")
}

func (u *User) Create(resultCh chan<- types.Result) {
	if status, err := HttpRequest(fmt.Sprintf("%s/create", userServicePath), u); err != nil {

		fmt.Println("User creation (preparing) failed:", u.UUID)
		resultCh <- types.Result{Success: false, Message: err.Error(), Status: status}
		return
	}

	fmt.Println("User creation (preparing) successfully:", u.UUID)
	resultCh <- types.Result{Success: true}
}

func (u *User) Rollback() {
	fmt.Println("User creation rollback:", u.UUID)
	HttpRequest(fmt.Sprintf("%s/rollback", userServicePath), u)
}

func (u *User) Commit(resultCh chan<- types.Result) {
	if status, err := HttpRequest(fmt.Sprintf("%s/commit", userServicePath), u); err != nil {

		fmt.Println("User creation (commit) failed:", u.UUID)
		resultCh <- types.Result{Success: false, Message: err.Error(), Status: status}
		return
	}

	fmt.Println("User creation (commit) successfully:", u.UUID)
	resultCh <- types.Result{Success: true}
}
