package routers

import (
	"Twitta/global"
	"github.com/gin-gonic/gin"

	ginprometheus "github.com/zinclabs/go-gin-prometheus"
)

// SetPrometheus sets up prometheus metrics for giin
func SetPrometheus(app *gin.Engine) {
	if !global.ServerConfig.PrometheusEnable {
		return
	}

	p := ginprometheus.NewPrometheus("Twitta", []*ginprometheus.Metric{})
	p.Use(app)
}
