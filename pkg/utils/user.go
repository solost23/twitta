package utils

import (
	"Twitta/pkg/models"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) *models.User {
	return c.Value("user").(*models.User)
}
