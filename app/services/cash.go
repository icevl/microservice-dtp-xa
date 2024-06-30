package services

import (
	"app/types"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Cash struct {
	UUID string `json:"uuid"`
}

var cashServicePath = loadCashServicePath()

func loadCashServicePath() string {
	godotenv.Load(".env")
	return os.Getenv("CASH_SERVICE")
}

func (u *Cash) Create(resultCh chan<- types.Result) {
	if status, err := HttpRequest(fmt.Sprintf("%s/create", cashServicePath), u); err != nil {

		fmt.Println("Cash creation (preparing) failed:", u.UUID)
		resultCh <- types.Result{Success: false, Message: err.Error(), Status: status}
		return
	}

	fmt.Println("Cash creation (preparing) successfully:", u.UUID)
	resultCh <- types.Result{Success: true}
}

func (u *Cash) Rollback() {
	fmt.Println("Cash creation rollback:", u.UUID)
	HttpRequest(fmt.Sprintf("%s/rollback", cashServicePath), u)
}

func (u *Cash) Commit(resultCh chan<- types.Result) {
	if status, err := HttpRequest(fmt.Sprintf("%s/commit", cashServicePath), u); err != nil {

		fmt.Println("Cash creation (commit) failed:", u.UUID)
		resultCh <- types.Result{Success: false, Message: err.Error(), Status: status}
		return
	}

	fmt.Println("Cash creation (commit) successfully:", u.UUID)
	resultCh <- types.Result{Success: true}
}
