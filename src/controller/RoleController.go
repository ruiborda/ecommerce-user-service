package controller

import (
	"net/http"

	dto "github.com/ruiborda/ecommerce-user-service/src/dto/common"
	"github.com/ruiborda/ecommerce-user-service/src/dto/role"
	"github.com/ruiborda/ecommerce-user-service/src/model"
	"github.com/ruiborda/ecommerce-user-service/src/service"
	"github.com/ruiborda/ecommerce-user-service/src/service/impl"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/src/swagger"
)

type RoleController struct {
	roleService service.RoleService
}

func NewRoleController() *RoleController {
	return &RoleController{
		roleService: impl.NewRoleServiceImpl(),
	}
}

var _ = swagger.Swagger().Path("/api/v1/roles").
	Post(func(operation openapi.Operation) {
		operation.Summary("Create a new role").
			OperationID("CreateRole").
			Tag("RoleController").
			Consume(mime.ApplicationJSON).
			Produces(mime.ApplicationJSON).
			BodyParameter(func(param openapi.Parameter) {
				param.Description("Role object that needs to be added to the system").
					Required(true).
					SchemaFromDTO(&role.CreateRoleRequest{})
			}).
			Security("BearerAuth")
	}).Doc()

// CreateRole handles the creation of a new role
func (roleController *RoleController) CreateRole(c *gin.Context) {
	var createRoleRequest = &role.CreateRoleRequest{}

	if err := c.BindJSON(createRoleRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar que hay permisos especificados
	if len(createRoleRequest.Permissions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere al menos un permiso para crear un rol"})
		return
	}

	// Obtener permisos usando FindPermissionsByIds
	permissions := model.FindPermissionsByIds(createRoleRequest.Permissions)

	// Validar que todos los permisos solicitados existen
	if len(*permissions) == 0 || len(*permissions) != len(createRoleRequest.Permissions) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Uno o más IDs de permisos no son válidos"})
		return
	}

	response := roleController.roleService.CreateRole(createRoleRequest)
	c.JSON(http.StatusCreated, response)
}

var _ = swagger.Swagger().Path("/api/v1/roles/{id}").
	Get(func(operation openapi.Operation) {
		operation.Summary("Get role by ID").
			OperationID("GetRoleByID").
			Tag("RoleController").
			Produces(mime.ApplicationJSON).
			PathParameter("id", func(param openapi.Parameter) {
				param.Description("ID of role to return").
					Required(true).
					Type("string")
			}).
			Response(http.StatusOK, func(response openapi.Response) {
				response.Description("Role object").
					SchemaFromDTO(&role.GetRoleByIdResponse{})
			}).
			Security("BearerAuth")
	}).Doc()

// GetRoleByID handles retrieval of a role by its ID
func (roleController *RoleController) GetRoleByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role ID is required"})
		return
	}

	// Validar que el ID sea un UUID válido
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	response := roleController.roleService.GetRoleById(id)
	if response == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

var _ = swagger.Swagger().Path("/api/v1/roles/{id}").
	Delete(func(operation openapi.Operation) {
		operation.Summary("Delete a role").
			OperationID("DeleteRole").
			Tag("RoleController").
			Produces(mime.ApplicationJSON).
			PathParameter("id", func(param openapi.Parameter) {
				param.Description("ID of role to delete").
					Required(true).
					Type("string")
			}).
			Security("BearerAuth")
	}).
	Doc()

// DeleteRole handles deletion of a role by its ID
func (roleController *RoleController) DeleteRole(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role ID is required"})
		return
	}

	// Validar que el ID sea un UUID válido
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	response := roleController.roleService.DeleteRoleById(id)
	if response == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process delete operation"})
		return
	}

	if !response.Success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": response.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

var _ = swagger.Swagger().Path("/api/v1/roles").
	Put(func(operation openapi.Operation) {
		operation.Summary("Update an existing role").
			OperationID("UpdateRole").
			Tag("RoleController").
			Consume(mime.ApplicationJSON).
			Produces(mime.ApplicationJSON).
			BodyParameter(func(param openapi.Parameter) {
				param.Description("Role object that needs to be updated").
					Required(true).
					SchemaFromDTO(&role.UpdateRoleRequest{})
			}).
			Security("BearerAuth")
	}).
	Doc()

// UpdateRole handles updating an existing role
func (roleController *RoleController) UpdateRole(c *gin.Context) {
	var updateRoleRequest = &role.UpdateRoleRequest{}

	if err := c.BindJSON(updateRoleRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateRoleRequest.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role ID is required"})
		return
	}

	// Validar que el ID sea un UUID válido
	if _, err := uuid.Parse(updateRoleRequest.Id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// Validar que hay permisos especificados
	if len(updateRoleRequest.Permissions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere al menos un permiso para actualizar un rol"})
		return
	}

	// Obtener permisos usando FindPermissionsByIds
	permissions := model.FindPermissionsByIds(updateRoleRequest.Permissions)

	// Validar que todos los permisos solicitados existen
	if len(*permissions) == 0 || len(*permissions) != len(updateRoleRequest.Permissions) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Uno o más IDs de permisos no son válidos"})
		return
	}

	response := roleController.roleService.UpdateRoleById(updateRoleRequest)
	if response == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

var _ = swagger.Swagger().Path("/api/v1/roles/pages").
	Get(func(operation openapi.Operation) {
		operation.Summary("Get roles with pagination").
			OperationID("GetAllByPageAndSize").
			Tag("RoleController").
			Produces(mime.ApplicationJSON).
			QueryParameter("page", func(param openapi.Parameter) {
				param.Description("Page number").
					Required(true).
					Type("integer")
			}).
			QueryParameter("size", func(param openapi.Parameter) {
				param.Description("Number of items per page").
					Required(true).
					Type("integer")
			}).
			QueryParameter("query", func(param openapi.Parameter) {
				param.Description("Search query").
					Required(false).
					Type("string")
			}).
			Security("BearerAuth")
	}).
	Doc()

// GetAllByPageAndSize handles retrieval of paginated roles
func (roleController *RoleController) GetAllByPageAndSize(c *gin.Context) {
	// Crear objeto Pageable desde parámetros de consulta
	pageable := dto.NewPageable(c.Query("page"), c.Query("size"), c.Query("query"))

	// Usar el nuevo método del servicio que maneja toda la paginación
	response := roleController.roleService.FindAllRolesPaginated(c, pageable)

	c.JSON(http.StatusOK, response)
}
