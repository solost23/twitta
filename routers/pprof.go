package routers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func SetPProf(app *gin.Engine) {
	pprof.Register(app)
}
