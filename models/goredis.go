package models

import (
	"context"

	"github.com/IBAX-io/go-explorer/conf"
)

var ctx = context.Background()

type RedisParams struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (rp *RedisParams) Set() error {
	return conf.GetRedisDbConn().Conn().Set(ctx, rp.Key, rp.Value, 0).Err()
}

func (rp *RedisParams) Get() error {
	val, err := conf.GetRedisDbConn().Conn().Get(ctx, rp.Key).Result()
	//if err != nil && err != redis.Nil {
	//	return err
	//}
	if err != nil {
		return err
