package service

import (
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
}
