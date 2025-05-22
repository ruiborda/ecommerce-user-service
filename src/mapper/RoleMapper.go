package mapper

import (
	"github.com/ruiborda/ecommerce-user-service/src/dto/role"
	"github.com/ruiborda/ecommerce-user-service/src/model"
)

type RoleMapper struct {
}

func (m *RoleMapper) CreateRoleRequestToRole(request *role.CreateRoleRequest) *model.Role {
	// Convertir IDs de permisos a objetos de permisos
	permissions := model.FindPermissionsByIds(request.Permissions)

	return &model.Role{
		Code:        request.Code,
		Permissions: permissions,
	}
}

func (m *RoleMapper) RoleToCreateRoleResponse(roleModel *model.Role) *role.CreateRoleResponse {
	permissions := roleModel.Permissions
	if permissions == nil {
		emptyPermissions := make([]model.Permission, 0)
		permissions = &emptyPermissions
	}

	return &role.CreateRoleResponse{
		Id:          roleModel.Id,
		Code:        roleModel.Code,
		Permissions: permissions,
	}
}

func (m *RoleMapper) RoleToGetRoleByIdResponse(roleModel *model.Role) *role.GetRoleByIdResponse {
	permissions := roleModel.Permissions
	if permissions == nil {
		emptyPermissions := make([]model.Permission, 0)
		permissions = &emptyPermissions
	}

	return &role.GetRoleByIdResponse{
		Id:          roleModel.Id,
		Code:        roleModel.Code,
		Permissions: permissions,
	}
}

func (m *RoleMapper) UpdateRoleRequestToRole(request *role.UpdateRoleRequest, existingModel *model.Role) *model.Role {
	permissions := model.FindPermissionsByIds(request.Permissions)

	existingModel.Code = request.Code
	existingModel.Permissions = permissions

	return existingModel
}

func (m *RoleMapper) RoleToUpdateRoleResponse(roleModel *model.Role) *role.UpdateRoleResponse {
	permissions := roleModel.Permissions
	if permissions == nil {
		emptyPermissions := make([]model.Permission, 0)
		permissions = &emptyPermissions
	}

	return &role.UpdateRoleResponse{
		Id:          roleModel.Id,
		Code:        roleModel.Code,
		Permissions: permissions,
	}
}

func (m *RoleMapper) RoleToDeleteRoleByIdResponse(roleId string, success bool) *role.DeleteRoleByIdResponse {
	return &role.DeleteRoleByIdResponse{
		Success: success,
		Message: getDeleteRoleMessage(roleId, success),
	}
}

func getDeleteRoleMessage(roleId string, success bool) string {
	if success {
		return "Role with ID " + roleId + " was successfully deleted"
	}
	return "Failed to delete role with ID " + roleId
}

func (m *RoleMapper) RolesToGetRolesResponse(models []*model.Role) []*role.GetRoleByIdResponse {
	var responses []*role.GetRoleByIdResponse

	for _, roleModel := range models {
		permissions := roleModel.Permissions
		if permissions == nil {
			emptyPermissions := make([]model.Permission, 0)
			permissions = &emptyPermissions
		}

		response := &role.GetRoleByIdResponse{
			Id:          roleModel.Id,
			Code:        roleModel.Code,
			Permissions: permissions,
		}

		responses = append(responses, response)
	}

	return responses
}

func (m *RoleMapper) GetRoleByIdResponseArrayToArray(dtos []*role.GetRoleByIdResponse) []*role.GetRoleByIdResponse {
	return dtos
}

func (m *RoleMapper) RolesToGetRolesByIdsResponse(models []*model.Role) *role.GetRolesByIdsResponse {
	roleResponses := m.RolesToGetRolesResponse(models)
	return &role.GetRolesByIdsResponse{
		Roles: roleResponses,
	}
}
