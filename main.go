package main

import (
	"github.com/gin-gonic/gin"
	router2 "github.com/ruiborda/ecommerce-user-service/src/router"
	"github.com/ruiborda/go-swagger-generator/src/middleware"
	"log/slog"
	"os"
)

func main() {

	router := gin.Default()

	router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json",
		UIPath:   "/",
	}))

	router2.ApiRouter(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	slog.Info("Starting server http://localhost:" + port)
	router.Run(":" + port)
}
