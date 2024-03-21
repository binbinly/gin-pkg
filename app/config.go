package app

import (
	"time"

	"github.com/spf13/viper"
)

var (
	Conf = &Config{}
)

func SetDefaultConf(v *viper.Viper) {
	v.SetDefault("Env", "local")
	v.SetDefault("Name", "gin-app")
	v.SetDefault("Mode", "debug")
	v.SetDefault("JwtTimeout", 864000)
	v.SetDefault("LogLevel", "debug")
	v.SetDefault("LogDir", "./logs/")
	v.SetDefault("Debug", true)
	v.SetDefault("HTTP", ServerConfig{
		Network:      "tcp",
		Addr:         ":9500",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	v.BindEnv("name")
	v.BindEnv("debug")
}

type Config struct {
	Env        string
	Name       string
	DfsUrl     string
	Mode       string
	LogLevel   string
	LogDir     string
	JwtSecret  string
	JwtTimeout int64
	Debug      bool
	Proxy      bool // 是否开启代理 http://[host]/ws -> ws://[host]
	HTTP       ServerConfig
	Websocket  ServerConfig
}

type ServerConfig struct {
	Network      string
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
