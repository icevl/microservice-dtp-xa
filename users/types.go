package main

type OperationPayload struct {
	UUID string `json:"uuid"`
}

type CreatePayload struct {
	UUID  string `json:"uuid"`
	Email string `json:"email"`
}

type ErrorResponse struct {
	Message string
	Code    int
}
