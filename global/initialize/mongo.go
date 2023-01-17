package initialize

import (
	"context"
	"time"
	"twitta/global"
	"twitta/pkg/mongoos"
)

func InitMongo() {
	if global.DB != nil {
		return
	}
	ctx := context.Background()
	var err error
	global.DB, err = mongoos.NewMongoConnect(ctx, &mongoos.Config{
		Hosts:      global.ServerConfig.MongoConfig.Hosts,
		AuthSource: global.ServerConfig.MongoConfig.AuthSource,
		Username:   global.ServerConfig.MongoConfig.Username,
		Password:   global.ServerConfig.MongoConfig.Password,
		Timeout:    time.Duration(global.ServerConfig.MongoConfig.Timeout),
	})
	if err != nil {
		panic(err)
	}
}
