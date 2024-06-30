package services

import (
	"app/types"
)

type Cash struct {
	UUID string
}

func (u *Cash) Create(resultCh chan<- types.Result) {
	resultCh <- types.Result{Success: true}
}

func (u *Cash) Rollback() {

}

func (u *Cash) Commit(resultCh chan<- types.Result) {
	resultCh <- types.Result{Success: true}
}
