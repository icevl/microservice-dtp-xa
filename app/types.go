package main

type NewUserPayload struct {
	Email string `json:"email" binding:"required"`
}
