package initialize

// 初始化全局对象
func Initialize() {
	// 初始化 Logger
	InitLogger()
	// 初始化配置
	InitConfig("config/config.yml")
	// 初始化location
	InitLoc()
	// 初始化 Mongo
	InitMongo()
}
