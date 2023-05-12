package domain

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/solost23/protopb/gen/go/protos/oss"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"twitta/global"
)

func NewOSSClient() oss.OssClient {
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
	return oss.NewOssClient(cc)
}
