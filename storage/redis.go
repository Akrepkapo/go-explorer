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
	rc = redis.NewClient(&redis.Options{
		Addr:     r.Str(),
		Password: r.Password,
		DB:       0,
	})
	_, err := rc.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisModel) Conn() *redis.Client {
	return rc
}
func (l *RedisModel) Close() error {
	return rc.Close()
}
