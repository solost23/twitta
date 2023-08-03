package initialize

import (
	"fmt"

	"twitta/global"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/solost23/protopb/gen/go/elastic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitElasticClient() {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port)

	target := fmt.Sprintf("consul://%s:%d/%s",
		global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port, global.ServerConfig.ESSrvConfig.Name)

	cc, err := grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		panic(err)
	}

	global.EsSrvClient = elastic.NewSearchServiceClient(cc)
}
