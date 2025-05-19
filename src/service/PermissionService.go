package service

import (
	"github.com/ruiborda/ecommerce-user-service/src/dto/permission"
)

type PermissionService interface {
	GetAllPermissions() *permission.GetAllPermissionsResponse
	GetPermissionById(id int) *permission.GetPermissionByIdResponse
	GetPermissionsByIds(request *permission.GetPermissionsByIdsRequest) permission.GetPermissionsByIdsResponse
}
