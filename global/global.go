package global

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/solost23/protopb/gen/go/elastic"
	"github.com/solost23/protopb/gen/go/oss"
	"github.com/solost23/protopb/gen/go/push"
	"github.com/solost23/protopb/gen/go/recognition"
	"twitta/configs"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	DB                       *mongo.Database
	Loc                      *time.Location
	ServerConfig             = &configs.ServerConfig{}
	RedisMapPool             = make(map[int]*redis.Client)
	EsSrvClient              elastic.SearchServiceClient
	PushSrvClient            push.PushServiceClient
	FaceRecognitionSrvClient recognition.FaceRecognitionServiceClient
	OssSrvClient             oss.OSSServiceClient
)
