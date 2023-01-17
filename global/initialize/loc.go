package initialize

import (
	"time"
	"twitta/global"
)

func InitLoc() {
	global.Loc, _ = time.LoadLocation(global.ServerConfig.TimeLocation)
}
