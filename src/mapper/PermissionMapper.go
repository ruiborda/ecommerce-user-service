package mapper

import (
	"github.com/ruiborda/ecommerce-user-service/src/dto/permission"
	"github.com/ruiborda/ecommerce-user-service/src/model"
)

type PermissionMapper struct {}

func (m *PermissionMapper) PermissionToGetPermissionByIdResponse(modelPermission *model.Permission) *permission.GetPermissionByIdResponse {
	return &permission.GetPermissionByIdResponse{
		Id:          modelPermission.Id,
		Method:      modelPermission.Method,
		Path:        modelPermission.Path,
		Name:        modelPermission.Name,
		Description: modelPermission.Description,
	}
}

func (m *PermissionMapper) PermissionsToGetAllPermissionsResponse(modelPermissions *[]model.Permission) *permission.GetAllPermissionsResponse {
	response := &permission.GetAllPermissionsResponse{
		Permissions: make([]permission.GetPermissionByIdResponse, 0, len(*modelPermissions)),
	}

	for _, perm := range *modelPermissions {
		permCopy := perm 
		response.Permissions = append(response.Permissions, *m.PermissionToGetPermissionByIdResponse(&permCopy))
	}

	return response
}

func (m *PermissionMapper) PermissionsToArray(modelPermissions *[]model.Permission) []permission.GetPermissionByIdResponse {
	response := make([]permission.GetPermissionByIdResponse, 0, len(*modelPermissions))

	for _, perm := range *modelPermissions {
		permCopy := perm 
		response = append(response, *m.PermissionToGetPermissionByIdResponse(&permCopy))
	}

	return response
}

func (m *PermissionMapper) PermissionsToGetPermissionsByIdsResponse(modelPermissions *[]model.Permission) permission.GetPermissionsByIdsResponse {
	response := make(permission.GetPermissionsByIdsResponse, 0, len(*modelPermissions))

	for _, perm := range *modelPermissions {
		permCopy := perm 
		response = append(response, *m.PermissionToGetPermissionByIdResponse(&permCopy))
	}

	return response
}
