package routers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"twitta/global"
	"twitta/pkg/utils"
)

const (
	timeout         = "5s"
	interval        = "30s"
	deregisterAfter = "10s"
)

func Run(client *api.Client, app *gin.Engine) {
	var err error

	serverConfig := global.ServerConfig
	ip := "127.0.0.1"

	if serverConfig.Mode != gin.DebugMode {
		ip, err = utils.GetInternalIP()
		if err != nil {
			zap.S().Panic("failed to get internal ip", err.Error())
		}
		port, err := utils.GetFreePort()
		if err != nil {
			zap.S().Panic("failed to get internal port", err.Error())
		}
		global.ServerConfig.Port = port
	}

	addr := fmt.Sprintf("%s:%d", ip, serverConfig.Port)
	// 生成检查对象
	check := &api.AgentServiceCheck{
		Interval:                       interval,
		Timeout:                        timeout,
		HTTP:                           fmt.Sprintf("http://%s", addr),
		Status:                         api.HealthPassing,
		DeregisterCriticalServiceAfter: deregisterAfter,
	}

	// 生成注册对象
	registration := &api.AgentServiceRegistration{
		ID:      addr,
		Name:    serverConfig.Name,
		Tags:    []string{serverConfig.Name, "web"},
		Address: ip,
		Port:    serverConfig.Port,
		Check:   check,
	}

	if err = client.Agent().ServiceRegister(registration); err != nil {
		zap.S().Panic("err register service", err.Error())
	}

	// Run the server
	server := &http.Server{
		Addr:         addr,
		Handler:      app,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	global.ServerConfig.Addr = addr
	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.S().Panic("server closed unexpect", err.Error())
		}
	}()
}
