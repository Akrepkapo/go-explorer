package storage

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var rc *redis.Client
var ctx = context.Background()

type RedisModel struct {
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

func (r *RedisModel) Str() string {
	return fmt.Sprintf("%s:%d", r.Address, r.Port)
}
}
func (l *RedisModel) Close() error {
	return rc.Close()
}
