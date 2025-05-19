package impl

import (
	"encoding/json"
	"errors"
	"github.com/ruiborda/ecommerce-user-service/src/dto/auth"
	"github.com/ruiborda/ecommerce-user-service/src/model"
	"github.com/ruiborda/ecommerce-user-service/src/repository"
	"github.com/ruiborda/ecommerce-user-service/src/repository/impl"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ruiborda/go-jwt/src/application/ports/input"
	"github.com/ruiborda/go-jwt/src/domain/entity"
	input2 "github.com/ruiborda/go-jwt/src/infrastructure/adapters/input"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	userRepository repository.UserRepository
	roleRepository repository.RoleRepository
}

func NewAuthServiceImpl() *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepository: impl.NewUserRepositoryImpl(),
		roleRepository: impl.NewRoleRepositoryImpl(),
	}
}

func (s *AuthServiceImpl) LoginWithGoogle(request *auth.LoginWithGoogleRequestDTO) (*auth.LoginWithAnyResponse, error) {
	// Validate input
	if request.AccessToken == "" {
		return nil, errors.New("access token is required")
	}

	// Get user info from Google
	googleUserInfo, err := s.getUserInfoFromGoogle(request.AccessToken)
	if err != nil {
		return nil, err
	}

	// Find or create user
	user, err := s.findOrCreateUserFromGoogle(googleUserInfo)
	if err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := s.generateJWTToken(user)
	if err != nil {
		return nil, err
	}

	// Return response
	return &auth.LoginWithAnyResponse{
		Id:           user.Id,
		FullName:     user.FullName,
		GivenName:    googleUserInfo.GivenName,
		FamilyName:   googleUserInfo.FamilyName,
		ProfileImage: googleUserInfo.Picture,
		Email:        googleUserInfo.Email,
		Jwt:          token,
	}, nil
}

func (s *AuthServiceImpl) LoginWithEmail(request *auth.LoginWithEmailRequestDTO) (*auth.LoginWithAnyResponse, error) {
	// Validate input
	if request.Email == "" || request.Password == "" {
		return nil, errors.New("email and password are required")
	}

	// Find user by email directly through repository
	user, err := s.userRepository.FindByEmail(request.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := s.generateJWTToken(user)
	if err != nil {
		return nil, err
	}

	// Return response
	return &auth.LoginWithAnyResponse{
		Id:           user.Id,
		FullName:     user.FullName,
		GivenName:    "",
		FamilyName:   "",
		ProfileImage: user.PictureUrl,
		Email:        user.Email,
		Jwt:          token,
	}, nil
}

// Helper methods

func (s *AuthServiceImpl) getUserInfoFromGoogle(accessToken string) (*auth.GoogleUserInfoResponse, error) {
	url := "https://openidconnect.googleapis.com/v1/userinfo"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, errors.New("failed to create request to Google API")
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		slog.Error("Error making request to Google API", "error", err)
		return nil, errors.New("failed to connect to Google API")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Error reading response body", "error", err)
		return nil, errors.New("failed to read response from Google API")
	}

	if res.StatusCode == http.StatusUnauthorized {
		slog.Error("Unauthorized access to Google API", "status", res.StatusCode)
		return nil, errors.New("unauthorized Google access token")
	}

	googleUserInfoResponse := &auth.GoogleUserInfoResponse{}
	err = json.Unmarshal(body, googleUserInfoResponse)
	if err != nil {
		slog.Error("Error unmarshalling Google response", "error", err)
		return nil, errors.New("failed to parse Google user info")
	}

	return googleUserInfoResponse, nil
}

func (s *AuthServiceImpl) findOrCreateUserFromGoogle(googleUserInfo *auth.GoogleUserInfoResponse) (*model.User, error) {
	// Check if user exists by email directly from repository
	user, err := s.userRepository.FindByEmail(googleUserInfo.Email)

	if err == nil && user != nil {
		// User exists, update with Google info
		user.FullName = googleUserInfo.Name
		user.PictureUrl = googleUserInfo.Picture
		user.UpdatedAt = time.Now()

		// Update user in database using repository
		user, err = s.userRepository.Update(user)
		if err != nil {
			slog.Warn("Failed to update user from Google login", "email", googleUserInfo.Email, "error", err)
			// Continue anyway as this is just an update
		}

		return user, nil
	} else {
		// User doesn't exist, create a new one with USER role
		// Get all roles from repository
		roles, err := s.roleRepository.FindAll()
		if err != nil {
			slog.Error("Failed to fetch roles", "error", err)
			return nil, errors.New("failed to fetch roles")
		}

		// Find USER role
		var userRoleId string
		for _, role := range roles {
			if role.Code == "USER" {
				userRoleId = role.Id
				break
			}
		}

		if userRoleId == "" {
			slog.Error("USER role not found in the database")
			return nil, errors.New("required role not found")
		}

		// Create new user
		now := time.Now()
		newUser := &model.User{
			Email:                  googleUserInfo.Email,
			FullName:               googleUserInfo.Name,
			PictureUrl:             googleUserInfo.Picture,
			CreatedAt:              now,
			UpdatedAt:              now,
			RoleIds:                []string{userRoleId},
			FavoriteNewsArticleIds: []string{},
		}

		// Save directly to repository
		newUser, err = s.userRepository.Create(newUser)
		if err != nil {
			slog.Error("Failed to create user", "error", err)
			return nil, errors.New("failed to create user")
		}

		return newUser, nil
	}
}

func (s *AuthServiceImpl) generateJWTToken(user *model.User) (string, error) {
	var roleCodes []string
	var permissionIds []int

	// Get roles directly from repository
	if len(user.RoleIds) > 0 {
		roles, err := s.roleRepository.FindByIds(user.RoleIds)
		if err != nil {
			slog.Error("Failed to fetch roles for user", "userId", user.Id, "error", err)
			// Continue with empty roles/permissions
		} else {
			// Extract role codes and permission IDs
			for _, role := range roles {
				roleCodes = append(roleCodes, role.Code)
				if role.Permissions != nil {
					for _, permission := range *role.Permissions {
						permissionIds = append(permissionIds, permission.Id)
					}
				}
			}
		}
	}

	// Create JWT token
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		slog.Error("JWT_SECRET environment variable is not set")
		return "", errors.New("JWT secret not configured")
	}

	inputPort := input.NewJWTHS256InputPort[*auth.JwtPrivateClaims]([]byte(jwtSecret))
	inputAdapter := input2.NewJwtInputAdapter[*auth.JwtPrivateClaims](inputPort)

	jwt, err := inputAdapter.CreateJwt(
		&entity.JOSEHeader{
			Algorithm: "HS256",
			Type:      "JWT",
		},
		&entity.JWTClaims[*auth.JwtPrivateClaims]{
			RegisteredClaims: &entity.RegisteredClaims{
				Issuer:         "ecommerce-user-service",
				Subject:        user.Id,
				ExpirationTime: time.Now().Add(time.Hour * 24).Unix(),
			},
			PrivateClaims: &auth.JwtPrivateClaims{
				Email:         user.Email,
				Roles:         roleCodes,
				PermissionIds: permissionIds,
			},
		})

	if err != nil {
		slog.Error("Error creating JWT", "error", err)
		return "", errors.New("failed to generate token")
	}

	return jwt.Token.GetToken(), nil
}
