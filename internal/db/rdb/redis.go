package rdb

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	Client *redis.Client
}

func New(ctx context.Context) (*RedisDB, error) {
	const t = 10
	timeout, cancel := context.WithTimeout(ctx, t*time.Second)
	defer cancel()
	cl := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err := cl.Ping(timeout).Err(); err != nil {
		return nil, errors.New("can't connect to redis: " + err.Error())
	}
	return &RedisDB{
		Client: cl,
	}, nil
}

func (r *RedisDB) ZAdd(ctx context.Context, setName string, id string, score float64) error {
	return r.Client.ZAdd(ctx, setName, redis.Z{
		Score:  score,
		Member: id,
	}).Err()
}

func (r *RedisDB) ZRevRank(ctx context.Context, setName string, id string) (int64, error) {
	return r.Client.ZRevRank(ctx, setName, id).Result()
}

func (r *RedisDB) ZRem(ctx context.Context, setName string, id string) error {
	return r.Client.ZRem(ctx, setName, id).Err()
}

func (r *RedisDB) ZRange(ctx context.Context, setName string, start int64, end int64) ([]redis.Z, error) {
	return r.Client.ZRangeWithScores(ctx, setName, start, end).Result()
}

func (r *RedisDB) HSet(ctx context.Context, setName string, values ...interface{}) (int64, error) {
	return r.Client.HSet(ctx, setName, values...).Result()
}

func (r *RedisDB) HIncryBy(ctx context.Context, setName string, field string, incr int64) (int64, error) {
	return r.Client.HIncrBy(ctx, setName, field, incr).Result()
}

func (r *RedisDB) HGet(ctx context.Context, setName string, field string) (string, error) {
	return r.Client.HGet(ctx, setName, field).Result()
}

func (r *RedisDB) HGetUint(ctx context.Context, setName string, field string) (uint64, error) {
	return r.Client.HGet(ctx, setName, field).Uint64()
}

func (r *RedisDB) HGetAll(ctx context.Context, setName string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, setName).Result()
}
