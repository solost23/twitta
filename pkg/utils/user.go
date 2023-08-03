package utils

import (
	"twitta/pkg/models"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) *models.User {
	return c.MustGet("user").(*models.User)
}
