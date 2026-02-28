package httputil

import (
	"github.com/gin-gonic/gin"
	apperrors "github.com/kidboy-man/8-level-desent/app/errors"
)

func ReturnSuccess(c *gin.Context, statusCode int, data interface{}) {
	if data == nil {
		c.Status(statusCode)
		return
	}
	c.JSON(statusCode, data)
}

func ReturnSuccessWithPagination(c *gin.Context, statusCode int, data interface{}, total, page, limit int) {
	response := gin.H{
		"data":  data,
		"total": total,
	}
	if page > 0 && limit > 0 {
		response["page"] = page
		response["limit"] = limit
	}
	c.JSON(statusCode, response)
}

func ReturnError(c *gin.Context, err error) {
	appErr := apperrors.FromError(err)
	c.JSON(appErr.HTTPStatus, gin.H{
		"error":   appErr.Code,
		"message": appErr.Message,
	})
}

func AbortWithError(c *gin.Context, err error) {
	appErr := apperrors.FromError(err)
	c.AbortWithStatusJSON(appErr.HTTPStatus, gin.H{
		"error":   appErr.Code,
		"message": appErr.Message,
	})
}
