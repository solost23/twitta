package configs

type ServerConfig struct {
	Name                  string              `mapstructure:"name"`
	DebugMode             string              `mapstructure:"debug_mode"`
	Addr                  string              `mapstructure:"addr"`
	Port                  int                 `mapstructure:"port"`
	TimeLocation          string              `mapstructure:"time_location"`
	PrometheusEnable      bool                `mapstructure:"prometheus_enable"`
	ConfigPath            string              `mapstructure:"config_path"`
	MongoConfig           MongoConf           `mapstructure:"mongo"`
	RedisConfig           RedisConf           `mapstructure:"redis"`
	JWTConfig             JWTConf             `mapstructure:"jwt"`
	ConsulConfig          ConsulConf          `mapstructure:"consul"`
	StaticOSS             StaticOSSConf       `mapstructure:"static-oss"`
	PushSrvConfig         PushSrvConf         `mapstructure:"push"`
	OSSSrvConfig          OSSSrvConf          `mapstructure:"oss"`
	ESSrvConfig           ESSrvConf           `mapstructure:"es"`
	FaceRecognitionConfig FaceRecognitionConf `mapstructure:"face_recognition"`
}

type MongoConf struct {
	Hosts      []string `mapstructure:"hosts"`
	AuthSource string   `mapstructure:"auth_source"`
	Username   string   `mapstructure:"username"`
	Password   string   `mapstructure:"password"`
	Timeout    int      `mapstructure:"timeout"`
}

type RedisConf struct {
	Addr string `mapstructure:"addr"`
}

type JWTConf struct {
	Key      string `mapstructure:"key"`
	Duration int64  `mapstructure:"duration"`
}

type ConsulConf struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type StaticOSSConf struct {
	Domain string `mapstructure:"domain"`
}

type PushSrvConf struct {
	Name string `mapstructure:"name"`
}

type OSSSrvConf struct {
	Name string `mapstructure:"name"`
}

type ESSrvConf struct {
	Name string `mapstructure:"name"`
}

type FaceRecognitionConf struct {
	Name string `mapstructure:"name"`
}
