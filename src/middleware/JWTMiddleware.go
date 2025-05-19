package middleware

import (
	"UserService/src/dto/auth"
	"github.com/gin-gonic/gin"
	"github.com/ruiborda/go-jwt/src/application/ports/input"
	"github.com/ruiborda/go-jwt/src/domain/entity"
	input2 "github.com/ruiborda/go-jwt/src/infrastructure/adapters/input"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// RequireJWT middleware checks if a valid JWT token is present
func RequireJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Extract JWT token from bearer format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}
		token := tokenParts[1]

		// Verify token
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			slog.Error("JWT_SECRET environment variable is not set")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		inputPort := input.NewJWTHS256InputPort[*auth.JwtPrivateClaims]([]byte(jwtSecret))
		inputAdapter := input2.NewJwtInputAdapter[*auth.JwtPrivateClaims](inputPort)

		err := inputAdapter.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Extract claims from the token for use in handlers
		jwt := entity.NewJwtFromToken[*auth.JwtPrivateClaims](token)
		if jwt == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		// Store claims in context for later use
		c.Set("jwtClaims", jwt.Claims)
		c.Next()
	}
}

// RequirePermission middleware checks if user has the required permission ID
func RequirePermission(permissionId int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First ensure JWT middleware has been run
		claimsValue, exists := c.Get("jwtClaims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No JWT claims found"})
			return
		}

		claims, ok := claimsValue.(*entity.JWTClaims[*auth.JwtPrivateClaims])
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid JWT claims format"})
			return
		}

		// Check if user has the required permission
		hasPermission := false
		if claims.PrivateClaims != nil && claims.PrivateClaims.PermissionIds != nil {
			for _, id := range claims.PrivateClaims.PermissionIds {
				if id == permissionId {
					hasPermission = true
					break
				}
			}
		}

		if !hasPermission {
			slog.Info("Access denied: missing required permission", "requiredPermission", permissionId, "email", claims.PrivateClaims.Email)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
			return
		}

		c.Next()
	}
}
