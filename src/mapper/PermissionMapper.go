package mapper

import (
	"github.com/ruiborda/ecommerce-user-service/src/dto/permission"
	"github.com/ruiborda/ecommerce-user-service/src/model"
)

type PermissionMapper struct {
}

// PermissionToGetPermissionByIdResponse convierte un modelo Permission a un GetPermissionByIdResponse
func (m *PermissionMapper) PermissionToGetPermissionByIdResponse(modelPermission *model.Permission) *permission.GetPermissionByIdResponse {
	return &permission.GetPermissionByIdResponse{
		Id:          modelPermission.Id,
		Method:      modelPermission.Method,
		Path:        modelPermission.Path,
		Name:        modelPermission.Name,
		Description: modelPermission.Description,
	}
}

// PermissionsToGetAllPermissionsResponse convierte una lista de modelo Permissions a un GetAllPermissionsResponse
func (m *PermissionMapper) PermissionsToGetAllPermissionsResponse(modelPermissions *[]model.Permission) *permission.GetAllPermissionsResponse {
	response := &permission.GetAllPermissionsResponse{
		Permissions: make([]permission.GetPermissionByIdResponse, 0, len(*modelPermissions)),
	}

	for _, perm := range *modelPermissions {
		permCopy := perm // Crear una copia para evitar problemas con el puntero
		response.Permissions = append(response.Permissions, *m.PermissionToGetPermissionByIdResponse(&permCopy))
	}

	return response
}

// PermissionsToArray convierte una lista de modelo Permissions a un array de GetPermissionByIdResponse
// Este método devuelve directamente el array sin envolverlo en un objeto
func (m *PermissionMapper) PermissionsToArray(modelPermissions *[]model.Permission) []permission.GetPermissionByIdResponse {
	response := make([]permission.GetPermissionByIdResponse, 0, len(*modelPermissions))

	for _, perm := range *modelPermissions {
		permCopy := perm // Crear una copia para evitar problemas con el puntero
		response = append(response, *m.PermissionToGetPermissionByIdResponse(&permCopy))
	}

	return response
}

// PermissionsToGetPermissionsByIdsResponse convierte una lista de modelo Permissions a un GetPermissionsByIdsResponse
func (m *PermissionMapper) PermissionsToGetPermissionsByIdsResponse(modelPermissions *[]model.Permission) permission.GetPermissionsByIdsResponse {
	response := make(permission.GetPermissionsByIdsResponse, 0, len(*modelPermissions))

	for _, perm := range *modelPermissions {
		permCopy := perm // Crear una copia para evitar problemas con el puntero
		response = append(response, *m.PermissionToGetPermissionByIdResponse(&permCopy))
	}

	return response
}
