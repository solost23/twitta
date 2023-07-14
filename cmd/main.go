package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
	"go.uber.org/zap"
	"twitta/global"
	"twitta/global/initialize"
	"twitta/routers"
)

// @title twitta API
// @version 1.0.0
// @description twitta API documents

// @contact.name twitta

// @securityDefinitions.basic BasicAuth
// @host localhost:6565
// @BasePath /
// @schemes http https

var (
	version  = "__BUILD_VERSION_"
	execDir  string
	st, v, V bool
)

func init() {
	flag.StringVar(&execDir, "d", ".", "项目目录")
	flag.BoolVar(&v, "v", false, "查看版本号")
	flag.BoolVar(&V, "V", false, "查看版本号")
	flag.BoolVar(&st, "s", false, "项目状态")
	flag.Parse()
}

func main() {
	if v || V {
		fmt.Println(version)
		os.Exit(-1)
	}

	initialize.Initialize(execDir)
	// 初始化所需要的服务
	initialize.InitESClient()
	initialize.InitFaceRecognitionClient()
	initialize.InitOSSClient()
	initialize.InitPushClient()

	client, err := api.NewClient(&api.Config{
		Address:   fmt.Sprintf("%s:%d", global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port),
		Scheme:    "http",
		Transport: cleanhttp.DefaultTransport(),
	})
	if err != nil {
		zap.S().Panic(err)
	}

	// HTTP init
	app := gin.Default()
	routers.Setup(app)

	routers.Run(client, app)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-c
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
			// 注销consul服务
			if err = client.Agent().ServiceDeregister(global.ServerConfig.Addr); err != nil {
				zap.S().Error("consul unregister service failed", global.ServerConfig.Addr, err)
			}
			zap.S().Info("consul unregister service success")
			return
		default:
			return
		}
	}
}
