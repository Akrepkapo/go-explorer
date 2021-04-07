package conf

import (
	"fmt"
	"time"
)

type serverModel struct {
	Mode                 string        `yaml:"mode"`                    // run mode
	Host                 string        `yaml:"host"`                    // server host
	Port                 int           `yaml:"port"`                    // server port
	EnableHttps          bool          `yaml:"enable_https"`            // enable https
	CertFile             string        `yaml:"cert_file"`               // cert file path
	KeyFile              string        `yaml:"key_file"`                // key file path
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
