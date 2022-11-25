package global

import (
	"Twitta/configs"
	"github.com/solost23/go_interface/gen_go/oss"
	"github.com/solost23/go_interface/gen_go/push"
	"time"

	"github.com/go-redis/redis/v8"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	DB            *mongo.Client
	Loc           *time.Location
	ServerConfig  = &configs.ServerConfig{}
	RedisMapPool  = make(map[int]*redis.Client)
	PushSrvClient push.PushClient
	OSSSrvClient  oss.OssClient
)
