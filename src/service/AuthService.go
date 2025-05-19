package service

import (
	"github.com/ruiborda/ecommerce-user-service/src/dto/auth"
)

type AuthService interface {
	// LoginWithGoogle handles the Google OAuth login process
	LoginWithGoogle(request *auth.LoginWithGoogleRequestDTO) (*auth.LoginWithAnyResponse, error)

	// LoginWithEmail handles email/password authentication
	LoginWithEmail(request *auth.LoginWithEmailRequestDTO) (*auth.LoginWithAnyResponse, error)
}
