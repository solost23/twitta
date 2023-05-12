package global

import (
	"time"

	"twitta/configs"

	"github.com/go-redis/redis/v8"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	DB           *mongo.Client
	Loc          *time.Location
	ServerConfig = &configs.ServerConfig{}
	RedisMapPool = make(map[int]*redis.Client)
)
