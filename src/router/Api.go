package router

import (
	"github.com/ruiborda/ecommerce-user-service/src/controller"

	"github.com/gin-gonic/gin"
)

func ApiRouter(router *gin.Engine) {
	userController := controller.NewUserController()
	authController := controller.NewAuthController()
	roleController := controller.NewRoleController()

	v1 := router.Group("/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login/google", authController.LoginWithGoogle)
			auth.POST("/login/email", authController.LoginWithEmail)
		}

		// User routes
		users := v1.Group("/users")
		{
			users.POST("", userController.CreateUser)

			users.GET("", userController.GetAllUsers)

			users.GET("/:id", userController.GetUserById)

			users.GET("/email/:email", userController.GetUserByEmail)

			users.PUT("", userController.UpdateUserById)

			users.DELETE("/:id", userController.DeleteUserById)

			users.GET("/pages", userController.FindAllUsersByPageAndSize)

			users.POST("/by-ids", userController.GetUsersByIds)
		}

		// Role routes
		roles := v1.Group("/roles")
		{
			roles.POST("", roleController.CreateRole)

			roles.GET("/:id", roleController.GetRoleByID)

			roles.PUT("", roleController.UpdateRole)

			roles.DELETE("/:id", roleController.DeleteRole)

			roles.GET("/pages", roleController.GetAllByPageAndSize)
		}
	}
}
