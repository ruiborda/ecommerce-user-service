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

	// Auth routes - these should not be protected as they're for login
	router.POST(
		"/api/v1/auth/login/google",
		authController.LoginWithGoogle,
	)
	router.POST(
		"/api/v1/auth/login/email",
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
}
