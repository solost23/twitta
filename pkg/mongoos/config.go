package mongoos

import "time"

type Config struct {
	Hosts      []string
	AuthSource string
	Username   string
	Password   string
	Timeout    time.Duration // second
}
