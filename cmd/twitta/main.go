package main

import (
	"Twitta/global"
	"Twitta/global/initialize"
	"Twitta/routers"
	"Twitta/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Twitta API
// @version 1.0.0
// @description Twitta API documents

// @contact.name Twitta

// @securityDefinitions.basic BasicAuth
// @host localhost:6565
// @BasePath /
// @schemes http
func main() {
	initialize.Initialize()
	// Version
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("Twitta version: %s\n", global.ServerConfig.Version)
		os.Exit(0)
	}
	zap.S().Infof("Starting Twitta %s", global.ServerConfig.Version)

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
