package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/solost23/go_interface/gen_go/oss"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"twitta/global"
)

func InitOSSClient() {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port)

	target := fmt.Sprintf("consul://%s:%d/%s",
		global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port, global.ServerConfig.OSSSrvConfig.Name)

	cc, err := grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		panic(err)
	}

	global.OSSSrvClient = oss.NewOssClient(cc)
}
