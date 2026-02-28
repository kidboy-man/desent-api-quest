package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	httputil "github.com/kidboy-man/8-level-desent/app/controllers/http"
	apperrors "github.com/kidboy-man/8-level-desent/app/errors"
	"github.com/kidboy-man/8-level-desent/app/services"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			httputil.AbortWithError(c, apperrors.NewUnauthorized("missing authorization header"))
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			httputil.AbortWithError(c, apperrors.NewUnauthorized("invalid authorization header format"))
			return
		}

		claims, err := authService.ValidateToken(parts[1])
		if err != nil {
			httputil.AbortWithError(c, err)
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
