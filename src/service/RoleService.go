// filepath: /home/rui/ecommerce/UserService/src/service/RoleService.go
package service

import (
	"github.com/gin-gonic/gin"
	dto "github.com/ruiborda/ecommerce-user-service/src/dto/common"
	"github.com/ruiborda/ecommerce-user-service/src/dto/role"
)

type RoleService interface {
	CreateRole(request *role.CreateRoleRequest) *role.CreateRoleResponse
	GetRoleById(id string) *role.GetRoleByIdResponse
	GetAllRoles() []*role.GetRoleByIdResponse
	UpdateRoleById(request *role.UpdateRoleRequest) *role.UpdateRoleResponse
	DeleteRoleById(id string) *role.DeleteRoleByIdResponse
	FindAllRolesByPageAndSize(page, size int) []*role.GetRoleByIdResponse
	CountAllRoles() int64
	GetRolesByIds(ids []string) []*role.GetRoleByIdResponse
	// Nuevo método que maneja la paginación completa
	FindAllRolesPaginated(c *gin.Context, pageable *dto.Pageable) *dto.PaginationResponse[role.GetRoleByIdResponse]
}
