package routers

import (
	"github.com/gin-gonic/gin"
	ginPrometheus "github.com/zsais/go-gin-prometheus"
	"twitta/global"
)

// SetPrometheus sets up prometheus metrics for gin
func SetPrometheus(app *gin.Engine) {
	if !global.ServerConfig.PrometheusEnable {
		return
	}

	ginPrometheus.NewPrometheus(global.ServerConfig.Name).Use(app)
}
