package middlewares

import (
	"github.com/gin-gonic/gin"
	gcasbin "github.com/maxwellhertz/gin-casbin"
	"twitta/pkg/dao"
)

var casbinMiddleWare *gcasbin.CasbinMiddleware

func NewCasbinMiddleware() *gcasbin.CasbinMiddleware {
	if casbinMiddleWare != nil {
		return casbinMiddleWare
	}
	auth, err := gcasbin.NewCasbinMiddleware("configs/rbac_model.conf", "configs/rbac_policy.csv", func(c *gin.Context) string {
		user, exists := c.Get("user")
		if !exists {
			return ""
		}
		return user.(*dao.User).Role
	})
	if err != nil {
		panic(err)
	}
	casbinMiddleWare = auth
	return casbinMiddleWare
}
