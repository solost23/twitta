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
	var err error
	global.DB, err = mongoos.NewMongoConnect(context.Background(), &mongoos.Config{
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
