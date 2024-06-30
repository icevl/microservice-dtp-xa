package services

type UserServiceCreateParams struct {
	UUID  string
	Email string
}

func CreateUser(params UserServiceCreateParams) (string, error) {
	return params.Email, nil
}
