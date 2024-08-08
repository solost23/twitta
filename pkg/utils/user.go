package utils

import (
	"twitta/pkg/dao"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) *dao.User {
	return c.MustGet("user").(*dao.User)
}
