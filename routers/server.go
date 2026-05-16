package routers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"twitta/global"
	"twitta/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
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

	if gin.Mode() != gin.DebugMode {
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

	// 容器环境下通过 BIND_IP=0.0.0.0 让服务监听所有网卡
	bindIP := ip
	if envIP := os.Getenv("BIND_IP"); envIP != "" {
		bindIP = envIP
	}

	addr := fmt.Sprintf("%s:%d", ip, serverConfig.Port)
	bindAddr := fmt.Sprintf("%s:%d", bindIP, serverConfig.Port)

	check := &api.AgentServiceCheck{
		Interval:                       interval,
		Timeout:                        timeout,
		HTTP:                           fmt.Sprintf("http://%s", addr),
		Status:                         api.HealthPassing,
		DeregisterCriticalServiceAfter: deregisterAfter,
	}

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

	server := &http.Server{
		Addr:         bindAddr,
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
