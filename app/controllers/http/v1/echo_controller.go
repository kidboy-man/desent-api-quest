package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httputil "github.com/kidboy-man/8-level-desent/app/controllers/http"
)

type EchoController struct{}

func NewEchoController() *EchoController {
	return &EchoController{}
}

func (ctrl *EchoController) Echo(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		httputil.ReturnError(c, err)
		return
	}
	httputil.ReturnSuccess(c, http.StatusOK, body)
}
