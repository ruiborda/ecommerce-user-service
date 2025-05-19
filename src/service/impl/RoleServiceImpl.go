package impl

import (
	"UserService/src/dto/role"
	"UserService/src/mapper"
	"UserService/src/repository"
	"UserService/src/repository/impl"
	"log"
)

type RoleServiceImpl struct {
	roleRepository repository.RoleRepository
	roleMapper     *mapper.RoleMapper
}

func NewRoleServiceImpl() *RoleServiceImpl {
	return &RoleServiceImpl{
		roleRepository: impl.NewRoleRepositoryImpl(),
		roleMapper:     &mapper.RoleMapper{},
	}
}

// CreateRole crea un nuevo rol
func (s *RoleServiceImpl) CreateRole(request *role.CreateRoleRequest) *role.CreateRoleResponse {
	// Mapear request a modelo
	roleModel := s.roleMapper.CreateRoleRequestToRole(request)

	// Guardar en el repositorio
	createdRole, err := s.roleRepository.Create(roleModel)
	if err != nil {
		log.Printf("Error creating role: %v", err)
		return nil
	}

	// Mapear modelo a response
	return s.roleMapper.RoleToCreateRoleResponse(createdRole)
}

// GetRoleById obtiene un rol por su ID
func (s *RoleServiceImpl) GetRoleById(id string) *role.CreateRoleResponse {
	roleModel, err := s.roleRepository.FindById(id)
	if err != nil {
		log.Printf("Error getting role by ID: %v", err)
		return nil
	}

	if roleModel == nil {
		return nil
	}

	return s.roleMapper.RoleToGetRoleByIdResponse(roleModel)
}

// GetAllRoles obtiene todos los roles
func (s *RoleServiceImpl) GetAllRoles() []*role.CreateRoleResponse {
	roles, err := s.roleRepository.FindAll()
	if err != nil {
		log.Printf("Error getting all roles: %v", err)
		return nil
	}

	return s.roleMapper.RolesToGetRolesResponse(roles)
}

// UpdateRoleById actualiza un rol existente
func (s *RoleServiceImpl) UpdateRoleById(request *role.UpdateRoleRequest) *role.UpdateRoleResponse {
	// Primero obtener el rol existente
	existingRole, err := s.roleRepository.FindById(request.Id)
	if err != nil {
		log.Printf("Error finding role to update: %v", err)
		return nil
	}

	if existingRole == nil {
		log.Printf("Role not found with ID: %s", request.Id)
		return nil
	}

	// Actualizar el modelo de rol con datos de la solicitud
	updatedRoleModel := s.roleMapper.UpdateRoleRequestToRole(request, existingRole)

	// Guardar el rol actualizado
	savedRole, err := s.roleRepository.Update(updatedRoleModel)
	if err != nil {
		log.Printf("Error updating role: %v", err)
		return nil
	}

	return s.roleMapper.RoleToUpdateRoleResponse(savedRole)
}

// DeleteRoleById elimina un rol por su ID
func (s *RoleServiceImpl) DeleteRoleById(id string) *role.DeleteRoleByIdResponse {
	err := s.roleRepository.Delete(id)
	if err != nil {
		log.Printf("Error deleting role: %v", err)
		return nil
	}

	return s.roleMapper.RoleToDeleteRoleByIdResponse(id, true)
}
