package controller

import (
	"github.com/ruiborda/ecommerce-user-service/src/dto/user"
	"github.com/ruiborda/ecommerce-user-service/src/service"
	"github.com/ruiborda/ecommerce-user-service/src/service/impl"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/src/swagger"
)

type UserController struct {
	userService service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: impl.NewUserServiceImpl(),
	}
}

var _ = swagger.Swagger().Path("/v1/users").
	Post(func(operation openapi.Operation) {
		operation.Summary("Create a new user").
			OperationID("CreateUser").
			Tag("UserController").
			Consume(mime.ApplicationJSON).
			Produces(mime.ApplicationJSON).
			BodyParameter(func(param openapi.Parameter) {
				param.Description("User object that needs to be added to the system").
					Required(true).
					SchemaFromDTO(&user.CreateUserRequest{})
			})
	}).Doc()

func (userController *UserController) CreateUser(c *gin.Context) {
	var createUserRequest = &user.CreateUserRequest{}

	if err := c.BindJSON(createUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userController.userService.CreateUser(createUserRequest))
}

var _ = swagger.Swagger().Path("/v1/users/{id}").
	Get(func(operation openapi.Operation) {
		operation.Summary("Get user by ID").
			OperationID("GetUserById").
			Tag("UserController").
			Produces(mime.ApplicationJSON).
			PathParameter("id", func(param openapi.Parameter) {
				param.Description("ID of user to return").
					Required(true).
					Type("string")
			}).
			Response(http.StatusOK, func(response openapi.Response) {
				response.Description("User object").
					SchemaFromDTO(&user.GetUserByIdResponse{})
			})
	}).Doc()

func (userController *UserController) GetUserById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Validar que el ID sea un UUID válido
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	response := userController.userService.GetUserById(id)
	if response == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

var _ = swagger.Swagger().Path("/v1/users/email/{email}").
	Get(func(operation openapi.Operation) {
		operation.Summary("Get user by email").
			OperationID("GetUserByEmail").
			Tag("UserController").
			Produces(mime.ApplicationJSON).
			PathParameter("email", func(param openapi.Parameter) {
				param.Description("Email of user to return").
					Required(true).
					Type("string")
			})
	}).
	Doc()

func (userController *UserController) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	response := userController.userService.GetUserByEmail(email)
	if response == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

var _ = swagger.Swagger().Path("/v1/users").
	Get(func(operation openapi.Operation) {
		operation.Summary("Get all users").
			OperationID("GetAllUsers").
			Tag("UserController").
			Produces(mime.ApplicationJSON).
			Response(http.StatusOK, func(response openapi.Response) {
				response.Description("List of users").
					SchemaFromDTO(&[]*user.GetUserByIdResponse{})
			})
	}).Doc()

func (userController *UserController) GetAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, userController.userService.GetAllUsers())
}

var _ = swagger.Swagger().Path("/v1/users").
	Put(func(operation openapi.Operation) {
		operation.Summary("Update an existing user").
			OperationID("UpdateUserById").
			Tag("UserController").
			Consume(mime.ApplicationJSON).
			Produces(mime.ApplicationJSON).
			BodyParameter(func(param openapi.Parameter) {
				param.Description("User object that needs to be updated").
					Required(true).
					SchemaFromDTO(&user.UpdateUserRequest{})
			})
	}).
	Doc()

func (userController *UserController) UpdateUserById(c *gin.Context) {
	var updateUserRequest = &user.UpdateUserRequest{}

	if err := c.BindJSON(updateUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateUserRequest.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Validar que el ID sea un UUID válido
	if _, err := uuid.Parse(updateUserRequest.Id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	response := userController.userService.UpdateUserById(updateUserRequest)
	if response == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

var _ = swagger.Swagger().Path("/v1/users/{id}").
	Delete(func(operation openapi.Operation) {
		operation.Summary("Delete a user").
			OperationID("DeleteUserById").
			Tag("UserController").
			Produces(mime.ApplicationJSON).
			PathParameter("id", func(param openapi.Parameter) {
				param.Description("ID of user to delete").
					Required(true).
					Type("string")
			}).
			BodyParameter(func(param openapi.Parameter) {
				param.Description("User object that needs to be deleted").
					Required(true).
					SchemaFromDTO(&user.DeleteUserByIdResponse{})
			})
	}).
	Doc()

func (userController *UserController) DeleteUserById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Validar que el ID sea un UUID válido
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	response := userController.userService.DeleteUserById(id)
	if response == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process delete operation"})
		return
	}

	if !response.Success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": response.Message})
		return
	}

	c.JSON(http.StatusOK, response)
}

var _ = swagger.Swagger().Path("/v1/users/pages").
	Get(func(operation openapi.Operation) {
		operation.Summary("Get users with pagination").
			OperationID("FindAllUsersByPageAndSize").
			Tag("UserController").
			Produces(mime.ApplicationJSON).
			QueryParameter("page", func(param openapi.Parameter) {
				param.Description("Page number").
					Required(true).
					Type("integer")
			}).
			QueryParameter("size", func(param openapi.Parameter) {
				param.Description("Number of items per page").
					Required(true).
					Type("integer")
			})
	}).
	Doc()

func (userController *UserController) FindAllUsersByPageAndSize(c *gin.Context) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size parameter"})
		return
	}

	users := userController.userService.FindAllUsersByPageAndSize(page, size)

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"page":  page,
		"size":  size,
		"total": userController.userService.CountAllUsers(),
	})
}

var _ = swagger.Swagger().Path("/v1/users/by-ids").
	Post(func(operation openapi.Operation) {
		operation.Summary("Get users by IDs").
			OperationID("GetUsersByIds").
			Tag("UserController").
			Consume(mime.ApplicationJSON).
			Produces(mime.ApplicationJSON).
			QueryParameter("ids", func(p openapi.Parameter) {
				p.Description("Tags to filter by").
					Required(true).
					Type("array").
					CollectionFormat("multi").
					Items(func(item openapi.Schema) { item.Type("uuid") })
			}).
			Response(http.StatusOK, func(response openapi.Response) {
				response.Description("List of users").
					SchemaFromDTO(&[]*user.GetUserByIdResponse{})
			})
	}).
	Doc()

func (userController *UserController) GetUsersByIds(c *gin.Context) {
	ids := c.QueryArray("ids")

	if len(ids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDs are required"})
		return
	}

	users := userController.userService.GetUsersByIds(ids)
	if users == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
		return
	}

	c.JSON(http.StatusOK, users)
}
