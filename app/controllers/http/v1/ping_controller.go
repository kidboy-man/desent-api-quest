package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httputil "github.com/kidboy-man/8-level-desent/app/controllers/http"
)

type PingController struct{}

func NewPingController() *PingController {
	return &PingController{}
}

func (ctrl *PingController) Ping(c *gin.Context) {
	httputil.ReturnSuccess(c, http.StatusOK, gin.H{"message": "pong"})
}
