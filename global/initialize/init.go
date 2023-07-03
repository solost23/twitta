package initialize

import (
	"path"
)

const (
	WebConfigPath = "configs/config.yml"
	WebLogPath    = "logs"
)

// 初始化全局对象
func Initialize(execDir string) {
	// 初始化配置
	InitConfig(path.Join(execDir, WebConfigPath))
	// 初始化 Logger
	InitLogger(path.Join(execDir, WebLogPath))
	// 初始化location
	InitLoc()
	// 初始化 Mongo
	InitMongo()
}
