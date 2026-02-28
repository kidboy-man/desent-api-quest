package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	httputil "github.com/kidboy-man/8-level-desent/app/controllers/http"
	apperrors "github.com/kidboy-man/8-level-desent/app/errors"
)

type EchoController struct{}

func NewEchoController() *EchoController {
	return &EchoController{}
}

func (ctrl *EchoController) Echo(c *gin.Context) {
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		httputil.ReturnError(c, err)
		return
	}
	if !json.Valid(raw) {
		httputil.ReturnError(c, apperrors.NewBadRequest("invalid JSON"))
		return
	}
	c.Data(http.StatusOK, "application/json", raw)
}
