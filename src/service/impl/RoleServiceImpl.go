// filepath: /home/rui/ecommerce/UserService/src/service/impl/RoleServiceImpl.go
package impl

import (
	"log"

	"github.com/gin-gonic/gin"
	dto "github.com/ruiborda/ecommerce-user-service/src/dto/common"
	"github.com/ruiborda/ecommerce-user-service/src/dto/role"
	"github.com/ruiborda/ecommerce-user-service/src/mapper"
	"github.com/ruiborda/ecommerce-user-service/src/repository"
	"github.com/ruiborda/ecommerce-user-service/src/repository/impl"
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
func (s *RoleServiceImpl) GetRoleById(id string) *role.GetRoleByIdResponse {
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
func (s *RoleServiceImpl) GetAllRoles() []*role.GetRoleByIdResponse {
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

// FindAllRolesByPageAndSize obtiene roles paginados
func (s *RoleServiceImpl) FindAllRolesByPageAndSize(page, size int) []*role.GetRoleByIdResponse {
	roles, err := s.roleRepository.FindAllByPageAndSize(page, size)
	if err != nil {
		log.Printf("Error fetching paginated roles: %v", err)
		return nil
	}

	return s.roleMapper.RolesToGetRolesResponse(roles)
}

// CountAllRoles cuenta el número total de roles
func (s *RoleServiceImpl) CountAllRoles() int64 {
	count, err := s.roleRepository.Count()
	if err != nil {
		log.Printf("Error counting roles: %v", err)
		return 0
	}

	return count
}

// GetRolesByIds obtiene múltiples roles por sus IDs
func (s *RoleServiceImpl) GetRolesByIds(ids []string) []*role.GetRoleByIdResponse {
	roles, err := s.roleRepository.FindByIds(ids)
	if err != nil {
		log.Printf("Error fetching roles by IDs: %v", err)
		return nil
	}

	return s.roleMapper.RolesToGetRolesResponse(roles)
}

// FindAllRolesPaginated obtiene roles paginados y construye la respuesta paginada completa
func (s *RoleServiceImpl) FindAllRolesPaginated(c *gin.Context, pageable *dto.Pageable) *dto.PaginationResponse[role.GetRoleByIdResponse] {
	// Convert from one-based (client) to zero-based (service) pagination
	zeroBasedPage := pageable.Page - 1

	// Obtener datos de roles paginados
	roles, err := s.roleRepository.FindAllByPageAndSize(zeroBasedPage, pageable.Size)
	if err != nil {
		log.Printf("Error fetching paginated roles: %v", err)
		return nil
	}

	// Obtener el conteo total de roles
	totalElements, err := s.roleRepository.Count()
	if err != nil {
		log.Printf("Error counting roles: %v", err)
		return nil
	}

	// Convertir roles al formato de respuesta usando el mapper
	roleMapper := &mapper.RoleMapper{}
	rolesDTO := roleMapper.RolesToGetRolesResponse(roles)

	// Crear la respuesta paginada
	return dto.NewPaginationResponse(c, &rolesDTO, int(totalElements), pageable)
}
