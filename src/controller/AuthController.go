package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruiborda/ecommerce-user-service/src/dto/auth"
	"github.com/ruiborda/ecommerce-user-service/src/service"
	"github.com/ruiborda/ecommerce-user-service/src/service/impl"
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/src/swagger"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: impl.NewAuthServiceImpl(),
	}
}

var _ = swagger.Swagger().Path("/api/v1/login-with-google").
	Post(func(operation openapi.Operation) {
		operation.Summary("Login or signup with Google OAuth").
			OperationID("LoginWithGoogle").
			Tag("AuthController").
			Consume(mime.ApplicationJSON).
			Produces(mime.ApplicationJSON).
			BodyParameter(func(param openapi.Parameter) {
				param.Description("Google access token received from frontend OAuth flow").
					Required(true).
					SchemaFromDTO(&auth.LoginWithGoogleRequestDTO{})
			}).
			Response(http.StatusOK, func(response openapi.Response) {
				response.Description("Login response with user details and JWT token").
					SchemaFromDTO(&auth.LoginWithAnyResponse{})
			}).
			Security("BearerAuth")
	}).Doc()

func (authController *AuthController) LoginWithGoogle(c *gin.Context) {
	var loginRequest = &auth.LoginWithGoogleRequestDTO{}

	if err := c.BindJSON(loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := authController.authService.LoginWithGoogle(loginRequest)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

var _ = swagger.Swagger().Path("/api/v1/auth/login-with-email").
	Post(func(operation openapi.Operation) {
		operation.Summary("Login with email and password").
			OperationID("LoginWithEmail").
			Tag("AuthController").
			Consume(mime.ApplicationJSON).
			Produces(mime.ApplicationJSON).
			BodyParameter(func(param openapi.Parameter) {
				param.Description("Email and password credentials").
					Required(true).
					SchemaFromDTO(&auth.LoginWithEmailRequestDTO{})
			}).
			Response(http.StatusOK, func(response openapi.Response) {
				response.Description("Login response with user details and JWT token").
					SchemaFromDTO(&auth.LoginWithAnyResponse{})
			}).
			Security("BearerAuth")
	}).Doc()

func (authController *AuthController) LoginWithEmail(c *gin.Context) {
	var loginRequest = &auth.LoginWithEmailRequestDTO{}

	if err := c.BindJSON(loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := authController.authService.LoginWithEmail(loginRequest)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
