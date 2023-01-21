package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
	"twitta/global"
	"twitta/global/initialize"
	"twitta/routers"
	"twitta/services"
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
	WebConfigPath = "configs/config.yml"
	WebLogPath    = "logs"
	version       = "__BUILD_VERSION_"
	execDir       string
	st, v, V      bool
)

func main() {
	flag.StringVar(&execDir, "d", ".", "项目目录")
	flag.BoolVar(&v, "v", false, "查看版本号")
	flag.BoolVar(&V, "V", false, "查看版本号")
	flag.BoolVar(&st, "s", false, "项目状态")
	flag.Parse()
	if v || V {
		fmt.Println(version)
		os.Exit(-1)
	}

	initialize.Initialize(path.Join(execDir, WebConfigPath))
	// HTTP init
	app := gin.New()
	routers.Setup(app)

	// Run the server
	server := &http.Server{
		Addr:         global.ServerConfig.Addr,
		Handler:      app,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	// 启动消费函数
	go func() {
		services.ConsumeEmailMessage()
	}()
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			if err == http.ErrServerClosed {
				zap.S().Errorf("Server closed")
			} else {
				zap.S().Errorf("Server closed unexpect %s", err.Error())
			}
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-c
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			_ = server.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
