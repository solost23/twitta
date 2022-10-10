package routers

import (
	"Twitta/global"
	"Twitta/pkg/middlewares"
	"github.com/gin-gonic/gin"
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
