package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Api struct {
	Name              string        `default:"klik"`
	Host              string        `default:"127.0.0.1"`
	Port              string        `default:"8000"`
	ReadHeaderTimeout time.Duration `split_words:"true" default:"60s"`
	GracefulTimeout   time.Duration `split_words:"true" default:"8s"`

	RequestLog bool `split_words:"true" default:"false"`
	RunSwagger bool `split_words:"true" default:"true"`
}

func API() Api {
	var api Api
	envconfig.MustProcess("API", &api)
	return api
}
