package service

import (
	"github.com/ruiborda/ecommerce-user-service/src/dto/permission"
)

type PermissionService interface {
	GetAllPermissions() *permission.GetAllPermissionsResponse
	GetAllPermissionsAsArray() []permission.GetPermissionByIdResponse
	GetPermissionById(id int) *permission.GetPermissionByIdResponse
	GetPermissionsByIds(request *permission.GetPermissionsByIdsRequest) permission.GetPermissionsByIdsResponse
	GetPermissionsByIdsAsArray(request *permission.GetPermissionsByIdsRequest) []permission.GetPermissionByIdResponse
}
