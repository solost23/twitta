package initialize

import (
	"Twitta/global"
	"time"
)

func InitLoc() {
	global.Loc, _ = time.LoadLocation(global.ServerConfig.TimeLocation)
}
