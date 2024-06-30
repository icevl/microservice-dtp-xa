package types

type NewUserPayload struct {
	Email string `json:"email" binding:"required"`
}

type Result struct {
	Message string
	Success bool
	Status  int
}
