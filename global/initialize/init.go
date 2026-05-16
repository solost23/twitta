package initialize

import (
	"path"
	"twitta/pkg/cache"
	"twitta/pkg/ws"
)

const (
	WebConfigPath = "configs/config.yml"
	WebLogPath    = "logs"
)

// Initialize 初始化全局对象
func Initialize(execDir string) {
	// 初始化配置
	InitConfig(path.Join(execDir, WebConfigPath))
	// 初始化 Logger
	InitLogger(path.Join(execDir, WebLogPath))
	// 初始化location
	InitLoc()
	// 初始化 Mongo
	InitMongo()
	// 初始化 WebSocket Hub（依赖 Redis）
	rdb, err := cache.RedisConnFactory(cache.DefaultDB)
	if err != nil {
		panic("failed to connect redis for ws hub: " + err.Error())
	}
	ws.InitHub(rdb)
}
