package main

import (
	"Twitta/global"
	"Twitta/global/initialize"
	"Twitta/router"
	"Twitta/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	initialize.Initialize()
	gin.SetMode(global.ServerConfig.DebugMode)
	server := &http.Server{
		Addr:         global.ServerConfig.Addr,
		Handler:      router.InitRouter(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	// 启动消费函数
	go func() {
		services.ConsumeEmailMessage()
	}()
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s \n", err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-c
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			server.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
