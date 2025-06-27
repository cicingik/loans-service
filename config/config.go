// Package config ...
package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type (
	// AppConfig ...
	AppConfig struct {
		HTTPHost string `envconfig:"http_host"`
		Secret   string `envconfig:"secret"`
		DBConfig struct {
			DbDriver               string `envconfig:"driver"`
			DbPort                 string `envconfig:"port"`
			DbHost                 string `envconfig:"host"`
			DbName                 string `envconfig:"name"`
			DbUser                 string `envconfig:"user"`
			DbPassword             string `envconfig:"password"`
			DbDebug                int    `envconfig:"debug"`
			MaxConnLifetimeSeconds int    `envconfig:"max_conn_lifetime_seconds"`
			MaxOpenConns           int    `envconfig:"max_open_conns"`
			MaxIdleConns           int    `envconfig:"max_idle_conns"`
		} `envconfig:"db"`
		HTTPPort int  `envconfig:"http_port"`
		Debug    bool `envconfig:"debug"`
	}
)

var cfg *AppConfig

// LoadConfig ...
func LoadConfig() *AppConfig {
	var xfg AppConfig
	err := envconfig.Process(AppName, &xfg)
	if err != nil {
		panic(fmt.Sprintf("cannot read in config: %s", err))
	}

	cfg = &xfg
	return cfg
}
