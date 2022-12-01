package middlewares

import (
	"Twitta/global"
	"Twitta/pkg/constants"
	"Twitta/pkg/response"
	"Twitta/pkg/utils"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	mongoadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/gin-gonic/gin"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
)

func AuthCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := utils.GetUser(c).Role
		mongoConf := global.ServerConfig.MongoConfig
		mc := mongooptions.Client().ApplyURI(fmt.Sprintf("mongodb://%s", mongoConf.Hosts[0]))
		a, err := mongoadapter.NewAdapterWithCollectionName(mc, constants.Mongo, "casbin_rules")
		if err != nil {
			response.Error(c, 2001, err)
			return
		}
		e, err := casbin.NewEnforcer("./configs/rbac_model.conf", a)
		if err != nil {
			response.Error(c, 2001, err)
			return
		}
		err = e.LoadPolicy()
		if err != nil {
			response.Error(c, 2001, err)
			return
		}
		ok, err := e.Enforce(role, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			response.Error(c, 2001, err)
			return
		}
		if !ok {
			response.Error(c, 2001, errors.New("权限认证不通过"))
			return
		}
		c.Next()
	}
}
