package main

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HttpPort             uint16 `default:"80" envconfig:"HTTP_PORT"`
	HttpAddr             string `default:"0.0.0.0" envconfig:"HTTP_ADDR"`
	HttpConcurrency      int    `default:"1000" envconfig:"HTTP_CONCURRENCY"`
	HttpBodyLimit        int    `default:"2048" envconfig:"HTTP_BODY_LIMIT"`
	Prefork              bool   `default:"false" envconfig:"HTTP_PREFORK"`
	ProxyHeader          string `default:"" envconfig:"HTTP_PROXY_HEADER"`
	DB                   string `required:"true" envconfig:"DB_URL"`
	DBMaxConnections     int    `default:"50" envconfig:"DB_MAX_CONNECTIONS"`
	DBMaxIdleConnections int    `default:"2" envconfig:"DB_MAX_IDLE_CONNECTIONS"`
	NotificationsURL     string `required:"true" envconfig:"NOTIFICATIONS_URL"`
	NotificationsToken   string `required:"true" envconfig:"NOTIFICATIONS_TOKEN"`
}

func LoadAppConfig() Config {
	var config Config
	envconfig.MustProcess("", &config)
	return config
}
