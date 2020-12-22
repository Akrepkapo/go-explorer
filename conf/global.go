package conf

import (
	"fmt"
	"time"

type UrlModel struct {
	URL string `yaml:"base_url"`
}

type LogConfig struct {
	LogTo     string
	LogLevel  string
	LogFormat string
}

func (r *serverModel) Str() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
