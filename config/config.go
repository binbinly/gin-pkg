package config

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var conf *Config

// Config conf struct.
type Config struct {
	envPrefix string
	dir       string
	fileType  string //yaml, json, toml, default is yaml
	val       map[string]*viper.Viper
	mu        sync.Mutex
}

// New create a config instance.
func New(opts ...Option) *Config {
	c := Config{
		envPrefix: "app",
		fileType:  fileTypeYaml,
		val:       make(map[string]*viper.Viper),
	}
	for _, opt := range opts {
		opt(&c)
	}

	conf = &c

	return &c
}

// Load alias for config func.
func Load(filename string, val any, hook func(v *viper.Viper)) error {
	return conf.Load(filename, val, hook)
}

// LoadWithType load conf by file type.
func LoadWithType(filename string, hook func(v *viper.Viper)) (*viper.Viper, error) {
	return conf.LoadWithType(filename, hook)
}

// Load scan data to struct.
func (c *Config) Load(filename string, val any, hook func(v *viper.Viper)) error {
	v, err := c.LoadWithType(filename, hook)
	if err != nil {
		return err
	}

	if err = v.Unmarshal(&val); err != nil {
		return err
	}

	// 注册每次配置文件发生变更后都会调用的回调函数
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
		// 每次配置文件发生变化，需要重新将其反序列化到结构体中
		if err := v.Unmarshal(&val); err != nil {
			panic(fmt.Errorf("unmarshal config error: %s \n", err.Error()))
		}
	})

	// 监控配置文件变化
	v.WatchConfig()

	return nil
}

// LoadWithType load conf by file type.
func (c *Config) LoadWithType(filename string, hook func(v *viper.Viper)) (v *viper.Viper, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.val[filename]
	if ok {
		return v, nil
	}

	v, err = c.load(filename, hook)
	if err != nil {
		return nil, err
	}
	c.val[filename] = v
	return v, nil
}

// Load file.
func (c *Config) load(filename string, hook func(v *viper.Viper)) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(c.dir)
	v.SetConfigName(filename)
	v.SetConfigType(c.fileType)
	v.AutomaticEnv()
	v.SetEnvPrefix(c.envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if hook != nil {
		hook(v)
	}
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	log.Println("Using config file: ", v.ConfigFileUsed(), " settings: ", v.AllSettings())

	return v, nil
}
