package routers

import (
	"github.com/gin-gonic/gin"
	"twitta/global"
	"twitta/pkg/middlewares"
)

func Setup(app *gin.Engine) {
	gin.SetMode(global.ServerConfig.DebugMode)

	app.Use(middlewares.RequestLog())

	// Debug for gin
	if gin.Mode() == gin.DebugMode {
		SetPProf(app)
	}
	SetPrometheus(app) // Set up Prometheus.
	SetRouters(app)    // Set up all API routers.
}
