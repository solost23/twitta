package global

import (
	es_service "github.com/solost23/protopb/gen/go/protos/es"
	"github.com/solost23/protopb/gen/go/protos/oss"
	"github.com/solost23/protopb/gen/go/protos/push"
	"time"
	"twitta/configs"

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
	ESSrvClient   es_service.SearchClient
)
