package auth

import (
	"strings"

	appErrors "fullstack-backend/internal/errors"
	"fullstack-backend/internal/usecase"
	"fullstack-backend/utils"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware(authUc usecase.IAuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendJSONErrorResponse(c, appErrors.New(appErrors.ErrUnauthorized, "Missing Authorization header"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.SendJSONErrorResponse(c, appErrors.New(appErrors.ErrUnauthorized, "Invalid Authorization format. Expected 'Bearer <token>'"))
			return
		}

		tokenString := parts[1]
		claims, err := authUc.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			utils.SendJSONErrorResponse(c, err)
			return
		}

		// Store claims in gin Context
		c.Set("user_claims", claims)
		c.Set("jwt_token", tokenString)
		c.Next()
	}
}
