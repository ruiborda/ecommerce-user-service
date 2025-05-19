// filepath: /home/rui/ecommerce/UserService/src/mapper/RoleMapper.go
package mapper

import (
	"github.com/ruiborda/ecommerce-user-service/src/dto/role"
	"github.com/ruiborda/ecommerce-user-service/src/model"
)

type RoleMapper struct {
}

// CreateRoleRequestToRole convierte un CreateRoleRequest a un modelo Role
func (m *RoleMapper) CreateRoleRequestToRole(request *role.CreateRoleRequest) *model.Role {
	// Convertir IDs de permisos a objetos de permisos
	permissions := model.FindPermissionsByIds(request.Permissions)

	return &model.Role{
		Code:        request.Code,
		Permissions: permissions,
	}
}

// RoleToCreateRoleResponse convierte un modelo Role a un CreateRoleResponse
func (m *RoleMapper) RoleToCreateRoleResponse(model *model.Role) *role.CreateRoleResponse {
	return &role.CreateRoleResponse{
		Id:          model.Id,
		Code:        model.Code,
		Permissions: model.Permissions,
	}
}

// RoleToGetRoleByIdResponse convierte un modelo Role a un GetRoleByIdResponse
func (m *RoleMapper) RoleToGetRoleByIdResponse(model *model.Role) *role.GetRoleByIdResponse {
	return &role.GetRoleByIdResponse{
		Id:          model.Id,
		Code:        model.Code,
		Permissions: model.Permissions,
	}
}

// UpdateRoleRequestToRole actualiza un modelo Role existente con datos de UpdateRoleRequest
func (m *RoleMapper) UpdateRoleRequestToRole(request *role.UpdateRoleRequest, existingModel *model.Role) *model.Role {
	// Convertir IDs de permisos a objetos de permisos
	permissions := model.FindPermissionsByIds(request.Permissions)

	// Actualizar los campos existentes
	existingModel.Code = request.Code
	existingModel.Permissions = permissions

	return existingModel
}

// RoleToUpdateRoleResponse convierte un modelo Role a un UpdateRoleResponse
func (m *RoleMapper) RoleToUpdateRoleResponse(model *model.Role) *role.UpdateRoleResponse {
	return &role.UpdateRoleResponse{
		Id:          model.Id,
		Code:        model.Code,
		Permissions: model.Permissions,
	}
}

// RoleToDeleteRoleByIdResponse convierte un resultado de eliminación a un DeleteRoleByIdResponse
func (m *RoleMapper) RoleToDeleteRoleByIdResponse(roleId string, success bool) *role.DeleteRoleByIdResponse {
	return &role.DeleteRoleByIdResponse{
		Success: success,
		Message: getDeleteRoleMessage(roleId, success),
	}
}

// Función auxiliar para crear un mensaje adecuado para roles
func getDeleteRoleMessage(roleId string, success bool) string {
	if success {
		return "Role with ID " + roleId + " was successfully deleted"
	}
	return "Failed to delete role with ID " + roleId
}

// RolesToGetRolesResponse convierte una lista de modelos Role a una lista de GetRoleByIdResponse
func (m *RoleMapper) RolesToGetRolesResponse(models []*model.Role) []*role.GetRoleByIdResponse {
	var responses []*role.GetRoleByIdResponse

	for _, model := range models {
		response := &role.GetRoleByIdResponse{
			Id:          model.Id,
			Code:        model.Code,
			Permissions: model.Permissions,
		}

		responses = append(responses, response)
	}

	return responses
}

// RolesToGetRolesByIdsResponse convierte una lista de modelos Role a una respuesta GetRolesByIdsResponse
func (m *RoleMapper) RolesToGetRolesByIdsResponse(models []*model.Role) *role.GetRolesByIdsResponse {
	roleResponses := m.RolesToGetRolesResponse(models)
	return &role.GetRolesByIdsResponse{
		Roles: roleResponses,
	}
}
