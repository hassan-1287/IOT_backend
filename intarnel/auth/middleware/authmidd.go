package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)


func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing authoriziton header"})
			ctx.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invild authoriziton format"})
			ctx.Abort()
			return
		}
		tokenString := parts[1]
		claims, err := ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invild token"})
			ctx.Abort()
			return
		}
		ctx.Set("claims", claims)
		userID := claims["sub"]
		ctx.Set("ids", userID)
		ctx.Next()
	}
}
