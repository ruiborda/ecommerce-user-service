package service

import (
	"github.com/gin-gonic/gin"
	dto "github.com/ruiborda/ecommerce-user-service/src/dto/common"
	"github.com/ruiborda/ecommerce-user-service/src/dto/user"
)

type UserService interface {
	CreateUser(request *user.CreateUserRequest) *user.CreateUserResponse
	GetUserById(id string) *user.GetUserByIdResponse
	GetUserByEmail(email string) *user.GetUserByIdResponse
	GetAllUsers() []*user.GetUserByIdResponse
	UpdateUserById(request *user.UpdateUserRequest) *user.UpdateUserResponse
	DeleteUserById(id string) *user.DeleteUserByIdResponse
	FindAllUsersByPageAndSize(page, size int) []*user.GetUserByIdResponse
	CountAllUsers() int64
	GetUsersByIds(ids []string) []*user.GetUserByIdResponse
	// Nuevo método que maneja la paginación completa
	FindAllUsersPaginated(c *gin.Context, pageable *dto.Pageable) *dto.PaginationResponse[user.GetUserByIdResponse]
}
