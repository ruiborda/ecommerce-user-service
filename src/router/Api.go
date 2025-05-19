package router

import (
	"UserService/src/controller"

	"github.com/gin-gonic/gin"
)

func ApiRouter(router *gin.Engine) {
	userController := controller.NewUserController()
	v1 := router.Group("/v1")
	{
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
	}

}
