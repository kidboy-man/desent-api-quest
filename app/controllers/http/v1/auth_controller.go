package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httputil "github.com/kidboy-man/8-level-desent/app/controllers/http"
	"github.com/kidboy-man/8-level-desent/app/models"
	"github.com/kidboy-man/8-level-desent/app/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ctrl *AuthController) GenerateToken(c *gin.Context) {
	var req models.TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ReturnError(c, err)
		return
	}

	token, err := ctrl.authService.GenerateToken(req.Username, req.Password)
	if err != nil {
		httputil.ReturnError(c, err)
		return
	}

	httputil.ReturnSuccess(c, http.StatusOK, models.TokenResponse{Token: token})
}
