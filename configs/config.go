package configs

type ServerConfig struct {
	Version          string      `mapstructure:"version"`
	DebugMode        string      `mapstructure:"debug_mode"`
	TimeLocation     string      `mapstructure:"time_location"`
	Addr             string      `mapstructure:"addr"`
	PrometheusEnable bool        `mapstructure:"prometheus_enable"`
	MongoConfig      MongoConf   `mapstructure:"mongo"`
	RedisConfig      RedisConf   `mapstructure:"redis"`
	MinioConfig      MinioConf   `mapstructure:"minio"`
	JWTConfig        JWTConf     `mapstructure:"jwt"`
	Email            EmailConf   `mapstructure:"email"`
	Zinc             Zinc        `mapstructure:"zinc"`
	ConsulConfig     ConsulConf  `mapstructure:"consul"`
	PushSrvConfig    PushSrvConf `mapstructure:"push"`
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

type MinioConf struct {
	EndPoint        string `mapstructure:"end_point"`
	AccessKeyId     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	UserSsl         bool   `mapstructure:"user_ssl"`
}

type JWTConf struct {
	Key      string `mapstructure:"key"`
	Duration int64  `mapstructure:"duration"`
}

type EmailConf struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	Password       string `mapstructure:"password"`
	SendPersonName string `mapstructure:"send_person_name"`
	SendPersonAddr string `mapstructure:"send_person_addr"`
}

type Zinc struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type ConsulConf struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type PushSrvConf struct {
	Name string `mapstructure:"name"`
}
