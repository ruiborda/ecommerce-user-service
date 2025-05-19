package auth

type LoginWithEmailRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
