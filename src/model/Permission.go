package model

type Permission struct {
	Id          int    `json:"id" firestore:"id,omitempty"`
	Method      string `json:"method" firestore:"method,omitempty"`
	Path        string `json:"path" firestore:"path,omitempty"`
	Name        string `json:"name" firestore:"name,omitempty"`
	Description string `json:"description" firestore:"description,omitempty"`
}

var Permissions *[]Permission
var PermissionsMap *map[int]Permission

const (
	// Permission Management
	GetAllPermissions   = 301
	GetPermissionById   = 302
	GetPermissionsByIds = 303

	// Role Management
	CreateRole        = 401
	GetRoleById       = 402
	GetRolesPaginated = 404
	DeleteRole        = 405
	UpdateRole        = 406

	// User Management
	CreateUser        = 501
	GetUserById       = 502
	UpdateUser        = 503
	DeleteUser        = 504
	GetUsersPaginated = 505

	// Product Management
	CreateProduct        = 601
	GetProductById       = 602
	UpdateProduct        = 603
	DeleteProduct        = 604
	GetProductsPaginated = 605
	AdjustProductStock   = 606
	SearchProducts       = 607
)

func GetAllPermissionsMap() *map[int]Permission {
	if PermissionsMap != nil {
		return PermissionsMap
	}
	PermissionsMap = &map[int]Permission{
		GetAllPermissions: {
			Id:          GetAllPermissions,
			Method:      "GET",
			Path:        "/permissions",
			Name:        "Ver Todos los Permisos",
			Description: "Permiso para obtener todos los permisos del sistema",
		},
		GetPermissionById: {
			Id:          GetPermissionById,
			Method:      "GET",
			Path:        "/permissions/:id",
			Name:        "Ver Permiso por ID",
			Description: "Permiso para obtener un permiso específico por su ID",
		},
		GetPermissionsByIds: {
			Id:          GetPermissionsByIds,
			Method:      "POST",
			Path:        "/permissions/by-ids",
			Name:        "Ver Permisos por IDs",
			Description: "Permiso para obtener múltiples permisos por sus IDs",
		},
		CreateRole: {
			Id:          CreateRole,
			Method:      "POST",
			Path:        "/roles",
			Name:        "Crear Rol",
			Description: "Permiso para crear un nuevo rol con permisos asociados",
		},
		GetRoleById: {
			Id:          GetRoleById,
			Method:      "GET",
			Path:        "/roles/:id",
			Name:        "Ver Rol por ID",
			Description: "Permiso para obtener información detallada de un rol específico por su ID",
		},
		GetRolesPaginated: {
			Id:          GetRolesPaginated,
			Method:      "GET",
			Path:        "/roles/pages",
			Name:        "Ver Roles Paginados",
			Description: "Permiso para obtener roles de forma paginada para mejor rendimiento",
		},
		DeleteRole: {
			Id:          DeleteRole,
			Method:      "DELETE",
			Path:        "/roles/:id",
			Name:        "Eliminar Rol",
			Description: "Permiso para eliminar un rol específico del sistema por su ID",
		},
		UpdateRole: {
			Id:          UpdateRole,
			Method:      "PUT",
			Path:        "/roles",
			Name:        "Actualizar Rol",
			Description: "Permiso para actualizar la información y permisos de un rol existente",
		},
		CreateUser: {
			Id:          CreateUser,
			Method:      "POST",
			Path:        "/users",
			Name:        "Crear Usuario",
			Description: "Permiso para crear un nuevo usuario en el sistema",
		},
		GetUserById: {
			Id:          GetUserById,
			Method:      "GET",
			Path:        "/users/:id",
			Name:        "Ver Usuario por ID",
			Description: "Permiso para recuperar información de un usuario específico por su ID",
		},
		UpdateUser: {
			Id:          UpdateUser,
			Method:      "PUT",
			Path:        "/users",
			Name:        "Actualizar Usuario",
			Description: "Permiso para actualizar la información de un usuario existente",
		},
		DeleteUser: {
			Id:          DeleteUser,
			Method:      "DELETE",
			Path:        "/users/:id",
			Name:        "Eliminar Usuario",
			Description: "Permiso para eliminar un usuario del sistema",
		},
		GetUsersPaginated: {
			Id:          GetUsersPaginated,
			Method:      "GET",
			Path:        "/users/pages",
			Name:        "Ver Usuarios Paginados",
			Description: "Permiso para obtener usuarios de forma paginada",
		},
		CreateProduct: {
			Id:          CreateProduct,
			Method:      "POST",
			Path:        "/products",
			Name:        "Crear Producto",
			Description: "Permiso para crear un nuevo producto en el catálogo",
		},
		GetProductById: {
			Id:          GetProductById,
			Method:      "GET",
			Path:        "/products/:id",
			Name:        "Ver Producto por ID",
			Description: "Permiso para obtener los detalles de un producto específico por su ID",
		},
		UpdateProduct: {
			Id:          UpdateProduct,
			Method:      "PUT",
			Path:        "/products",
			Name:        "Actualizar Producto",
			Description: "Permiso para modificar la información de un producto existente",
		},
		DeleteProduct: {
			Id:          DeleteProduct,
			Method:      "DELETE",
			Path:        "/products/:id",
			Name:        "Eliminar Producto",
			Description: "Permiso para eliminar un producto del catálogo por su ID",
		},
		GetProductsPaginated: {
			Id:          GetProductsPaginated,
			Method:      "GET",
			Path:        "/products/pages",
			Name:        "Ver Productos Paginados",
			Description: "Permiso para obtener productos de forma paginada",
		},
		AdjustProductStock: {
			Id:          AdjustProductStock,
			Method:      "PATCH",
			Path:        "/products/:id/stock",
			Name:        "Ajustar Stock de Producto",
			Description: "Permiso para incrementar o decrementar el stock de un producto",
		},
		SearchProducts: {
			Id:          SearchProducts,
			Method:      "GET",
			Path:        "/products/search",
			Name:        "Buscar Productos con Filtros",
			Description: "Permiso para buscar productos por términos y filtros avanzados (similar a Amazon)",
		},
	}
	return PermissionsMap
}

func GetAllPermissionsAsSlice() *[]Permission {
	if Permissions != nil {
		return Permissions
	}
	permissions := GetAllPermissionsMap()
	var permissionsSlice []Permission
	for _, permission := range *permissions {
		permissionsSlice = append(permissionsSlice, permission)
	}
	Permissions = &permissionsSlice
	return Permissions
}

func FindPermissionById(id int) *Permission {
	permissions := GetAllPermissionsMap()
	if permission, ok := (*permissions)[id]; ok {
		return &permission
	}
	return nil
}

func FindPermissionsByIds(ids []int) *[]Permission {
	permissions := GetAllPermissionsMap()
	var foundPermissions []Permission
	for _, id := range ids {
		if permission, ok := (*permissions)[id]; ok {
			foundPermissions = append(foundPermissions, permission)
		}
	}
	return &foundPermissions
}
