package global

import (
	"Twitta/config"
	"time"

	"github.com/go-redis/redis/v8"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	DB           *mongo.Client
	Loc          *time.Location
	ServerConfig = &config.ServerConfig{}
	RedisMapPool = make(map[int]*redis.Client)
)
