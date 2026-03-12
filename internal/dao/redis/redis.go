package redis

import (
	"context"

	"MicroCraft/internal/config"
	perr "MicroCraft/pkg/errors"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() error {
	if config.Conf == nil {
		return perr.New(perr.CONFIG_NOT_INIT)
	}
	rc := config.Conf.Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     rc.Addr,
		Password: rc.Password,
		DB:       rc.DB,
	})
	if _, err := RDB.Ping(context.Background()).Result(); err != nil {
		return perr.Wrap(perr.REDIS_CONNECT_FAIL, err)
	}

	return nil
}