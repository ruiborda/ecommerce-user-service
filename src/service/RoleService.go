package service

import (
	"UserService/src/dto/role"
)

type RoleService interface {
	CreateRole(request *role.CreateRoleRequest) *role.CreateRoleResponse
	GetRoleById(id string) *role.CreateRoleResponse
	GetAllRoles() []*role.CreateRoleResponse
	UpdateRoleById(request *role.UpdateRoleRequest) *role.UpdateRoleResponse
	DeleteRoleById(id string) *role.DeleteRoleByIdResponse
}
