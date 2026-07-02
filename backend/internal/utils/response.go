package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/dto"
)

func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, dto.ErrorResponse{Error: message})
}

func ValidationError(c *gin.Context, details map[string]string) {
	c.JSON(http.StatusBadRequest, dto.ErrorResponse{
		Error:   "validation failed",
		Details: details,
	})
}

func NotFound(c *gin.Context, resource string) {
	Error(c, http.StatusNotFound, resource+" not found")
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}
