package impl

import (
	"github.com/ruiborda/ecommerce-user-service/src/dto/permission"
	"github.com/ruiborda/ecommerce-user-service/src/mapper"
	"github.com/ruiborda/ecommerce-user-service/src/model"
)

type PermissionServiceImpl struct {
	permissionMapper *mapper.PermissionMapper
}

func NewPermissionServiceImpl() *PermissionServiceImpl {
	return &PermissionServiceImpl{
		permissionMapper: &mapper.PermissionMapper{},
	}
}

// GetAllPermissions obtiene todos los permisos del sistema
func (s *PermissionServiceImpl) GetAllPermissions() *permission.GetAllPermissionsResponse {
	permissions := model.GetAllPermissionsAsSlice()
	return s.permissionMapper.PermissionsToGetAllPermissionsResponse(permissions)
}

// GetAllPermissionsAsArray obtiene todos los permisos del sistema como un array directamente
func (s *PermissionServiceImpl) GetAllPermissionsAsArray() []permission.GetPermissionByIdResponse {
	permissions := model.GetAllPermissionsAsSlice()
	return s.permissionMapper.PermissionsToArray(permissions)
}

// GetPermissionById obtiene un permiso por su ID
func (s *PermissionServiceImpl) GetPermissionById(id int) *permission.GetPermissionByIdResponse {
	permissionModel := model.FindPermissionById(id)
	if permissionModel == nil {
		return nil
	}
	return s.permissionMapper.PermissionToGetPermissionByIdResponse(permissionModel)
}

// GetPermissionsByIds obtiene múltiples permisos por sus IDs
func (s *PermissionServiceImpl) GetPermissionsByIds(request *permission.GetPermissionsByIdsRequest) permission.GetPermissionsByIdsResponse {
	permissions := model.FindPermissionsByIds(request.Ids)
	return s.permissionMapper.PermissionsToGetPermissionsByIdsResponse(permissions)
}

// GetPermissionsByIdsAsArray obtiene múltiples permisos por sus IDs como un array directamente
func (s *PermissionServiceImpl) GetPermissionsByIdsAsArray(request *permission.GetPermissionsByIdsRequest) []permission.GetPermissionByIdResponse {
	permissions := model.FindPermissionsByIds(request.Ids)
	return s.permissionMapper.PermissionsToArray(permissions)
}
