package router

import (
	"github.com/ruiborda/ecommerce-user-service/src/controller"
	"github.com/ruiborda/ecommerce-user-service/src/middleware"
	"github.com/ruiborda/ecommerce-user-service/src/model"

	"github.com/gin-gonic/gin"
)

func ApiRouter(router *gin.Engine) {
	userController := controller.NewUserController()
	authController := controller.NewAuthController()
	roleController := controller.NewRoleController()
	permissionController := controller.NewPermissionController()

	// Auth routes - these should not be protected as they're for login
	router.POST(
		"/api/v1/auth/login-with-google",
		authController.LoginWithGoogle,
	)
	router.POST(
		"/api/v1/auth/login-with-email",
		authController.LoginWithEmail,
	)

	// User routes - protected with JWT and specific permissions
	router.POST(
		"/api/v1/users",
		middleware.RequireJWT(),
		middleware.RequirePermission(model.CreateUser),
		userController.CreateUser,
	)

	router.GET(
		"/api/v1/users",
		middleware.RequireJWT(),
		userController.GetAllUsers,
	)

	router.GET(
		"/api/v1/users/:id",
		middleware.RequireJWT(),
		middleware.RequirePermission(model.GetUserById),
		userController.GetUserById,
	)

	router.PUT(
		"/api/v1/users",
		middleware.RequireJWT(),
		middleware.RequirePermission(model.UpdateUser),
		userController.UpdateUserById,
	)

	router.DELETE(
		"/api/v1/users/:id",
		middleware.RequireJWT(),
		middleware.RequirePermission(model.DeleteUser),
		userController.DeleteUserById,
	)

	router.GET(
		"/api/v1/users/pages",
		middleware.RequireJWT(),
		middleware.RequirePermission(model.GetUsersPaginated),
		userController.FindAllUsersByPageAndSize,
	)

	// Role routes - protected with JWT and specific permissions
	router.POST(
		"/api/v1/roles",
		middleware.RequireJWT(),
		middleware.RequirePermission(model.CreateRole),
		roleController.CreateRole,
	)

	router.GET(
		"/api/v1/roles/:id",
		middleware.RequireJWT(),
		middleware.RequirePermission(model.GetRoleById),
		roleController.GetRoleByID,
	)

	router.PUT(
		"/api/v1/roles",
		middleware.RequireJWT(),
		middleware.RequirePermission(model.UpdateRole),
		roleController.UpdateRole,
	)

	router.DELETE(
		"/api/v1/roles/:id",
		middleware.RequireJWT(),
		middleware.RequirePermission(model.DeleteRole),
		roleController.DeleteRole,
	)

	router.GET(
		"/api/v1/roles/pages",
		middleware.RequireJWT(),
		middleware.RequirePermission(model.GetRolesPaginated),
		roleController.GetAllByPageAndSize,
	)

	// Permission routes - protected with JWT and specific permissions
	// Solo rutas de consulta ya que los permisos están en hard code
	router.GET(
		"/api/v1/permissions",
		middleware.RequireJWT(),
		middleware.RequirePermission(model.GetAllPermissions),
		permissionController.GetAllPermissions,
	)

	router.GET(
		"/api/v1/permissions/:id",
		middleware.RequireJWT(),
		middleware.RequirePermission(model.GetPermissionById),
		permissionController.GetPermissionById,
	)

	router.POST(
		"/api/v1/permissions/by-ids",
		middleware.RequireJWT(),
		middleware.RequirePermission(model.GetPermissionsByIds),
		permissionController.GetPermissionsByIds,
	)
}
