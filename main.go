package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	router2 "github.com/ruiborda/ecommerce-user-service/src/route"
	"github.com/ruiborda/go-swagger-generator/src/middleware"
	"github.com/ruiborda/go-swagger-generator/src/openapi"
	"github.com/ruiborda/go-swagger-generator/src/swagger"
	"log/slog"
	"os"
	"time"
)

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.Use(middleware.SwaggerGin(middleware.SwaggerConfig{
		Enabled:  true,
		JSONPath: "/openapi.json",
		UIPath:   "/",
	}))

	swagger.Swagger().SecurityDefinition("BearerAuth", func(sd openapi.SecurityScheme) {
		sd.Type("apiKey").
			Name("Authorization").
			In("header")
	})

	router2.ApiRouter(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	slog.Info("Starting server http://localhost:" + port)
	router.Run(":" + port)
}
