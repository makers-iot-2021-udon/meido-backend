package redis

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

const Nil = redis.Nil

const (
	tokenKey      = "TOKEN_"
	tokenDuration = time.Hour * 24
)

func New(dsn string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "",
		DB:       0,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to ping redis server")
	}
	return client, nil
}

func SetToken(cli *redis.Client, token string, userId int) error {
	if err := cli.Set(tokenKey+token, userId, tokenDuration).Err(); err != nil {
		return errors.Wrapf(err, "failed to set value")
	}
	return nil
}

func GetIDByToken(cli *redis.Client, token string) (int, error) {
	v, err := cli.Get(tokenKey + token).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get id from redis by token")
	}
	id, err := strconv.Atoi(v)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to convert string to inte")
	}
	return id, nil
}
