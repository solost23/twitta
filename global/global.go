package global

import (
	"time"

	"github.com/go-redis/redis/v8"
	es_service "github.com/solost23/protopb/gen/go/protos/es"
	"github.com/solost23/protopb/gen/go/protos/face_recognition"
	"github.com/solost23/protopb/gen/go/protos/oss"
	"github.com/solost23/protopb/gen/go/protos/push"
	"twitta/configs"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	DB                       *mongo.Client
	Loc                      *time.Location
	ServerConfig             = &configs.ServerConfig{}
	RedisMapPool             = make(map[int]*redis.Client)
	EsSrvClient              es_service.SearchClient
	PushSrvClient            push.PushClient
	FaceRecognitionSrvClient face_recognition.FaceRecognitionClient
	OssSrvClient             oss.OssClient
)
