package conf

import (
	"fmt"
	"time"
	JwtPriKeyPath        string        `yaml:"jwt_private_key_path"`    // jwt private key path
	TokenExpireSecond    time.Duration `yaml:"token_expire_second"`     // token expire second
	SystemStaticFilePath string        `yaml:"system_static_file_path"` // system static file path
}

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
