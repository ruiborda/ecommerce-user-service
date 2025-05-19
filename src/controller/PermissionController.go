package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ruiborda/ecommerce-user-service/src/dto/permission"
	"github.com/ruiborda/ecommerce-user-service/src/service"
	"github.com/ruiborda/ecommerce-user-service/src/service/impl"
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/openapi_spec/mime"
	"github.com/ruiborda/go-swagger-generator/src/swagger"
)

type PermissionController struct {
	permissionService service.PermissionService
}

func NewPermissionController() *PermissionController {
	return &PermissionController{
		permissionService: impl.NewPermissionServiceImpl(),
	}
}

var _ = swagger.Swagger().Path("/api/v1/permissions").
	Get(func(operation openapi.Operation) {
		operation.Summary("Get all permissions").
			OperationID("GetAllPermissions").
			Tag("PermissionController").
			Produces(mime.ApplicationJSON).
			Response(http.StatusOK, func(response openapi.Response) {
				response.Description("List of all system permissions").
					SchemaFromDTO(&permission.GetAllPermissionsResponse{})
			})
	}).Doc()

func (p *PermissionController) GetAllPermissions(c *gin.Context) {
	response := p.permissionService.GetAllPermissions()
	c.JSON(http.StatusOK, response)
}

var _ = swagger.Swagger().Path("/api/v1/permissions/{id}").
	Get(func(operation openapi.Operation) {
		operation.Summary("Get permission by ID").
			OperationID("GetPermissionById").
			Tag("PermissionController").
			Produces(mime.ApplicationJSON).
			PathParameter("id", func(param openapi.Parameter) {
				param.Description("ID of permission to return").
					Required(true).
					Type("integer")
			}).
			Response(http.StatusOK, func(response openapi.Response) {
				response.Description("Permission object").
					SchemaFromDTO(&permission.GetPermissionByIdResponse{})
			}).
			Response(http.StatusBadRequest, func(response openapi.Response) {
				response.Description("Invalid ID format").
					SchemaFromDTO(&permission.ErrorResponse{})
			}).
			Response(http.StatusNotFound, func(response openapi.Response) {
				response.Description("Permission not found").
					SchemaFromDTO(&permission.ErrorResponse{})
			})
	}).Doc()

func (p *PermissionController) GetPermissionById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse := &permission.ErrorResponse{
			Success: false,
			Message: "ID inválido",
			Error:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	response := p.permissionService.GetPermissionById(id)
	if response == nil {
		errorResponse := &permission.ErrorResponse{
			Success: false,
			Message: "Permission no encontrado",
		}
		c.JSON(http.StatusNotFound, errorResponse)
		return
	}

	c.JSON(http.StatusOK, response)
}

var _ = swagger.Swagger().Path("/api/v1/permissions/by-ids").
	Post(func(operation openapi.Operation) {
		operation.Summary("Get permissions by IDs").
			OperationID("GetPermissionsByIds").
			Tag("PermissionController").
			Consume(mime.ApplicationJSON).
			Produces(mime.ApplicationJSON).
			BodyParameter(func(param openapi.Parameter) {
				param.Description("List of permission IDs to fetch").
					Required(true).
					SchemaFromDTO(&permission.GetPermissionsByIdsRequest{})
			}).
			Response(http.StatusOK, func(response openapi.Response) {
				response.Description("List of permissions that match the requested IDs").
					SchemaFromDTO(&permission.GetPermissionsByIdsResponse{})
			}).
			Response(http.StatusBadRequest, func(response openapi.Response) {
				response.Description("Invalid request format").
					SchemaFromDTO(&permission.ErrorResponse{})
			})
	}).Doc()

func (p *PermissionController) GetPermissionsByIds(c *gin.Context) {
	var request permission.GetPermissionsByIdsRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse := &permission.ErrorResponse{
			Success: false,
			Message: "Formato de solicitud inválido",
			Error:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	response := p.permissionService.GetPermissionsByIds(&request)
	c.JSON(http.StatusOK, response)
}
