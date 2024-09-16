package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type rdm struct {
	client *redis.Client
}

func (r *rdm) IsValidInterface() bool {
	return true
}

func (r *rdm) Set(ctx context.Context, duration time.Duration, key string, value string) error {
	return r.client.Set(ctx, key, value, duration).Err()
}

func (r *rdm) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *rdm) Get(ctx context.Context, key string) (string, error) {
	result := r.client.Get(ctx, key)
	if result.Err() != nil {
		zap.L().Error("Get error", zap.Error(result.Err()), zap.String("key", key))
		return "", result.Err()
	}

	return result.Val(), nil
}

func (r *rdm) GetKeys(ctx context.Context, mustInclude string) ([]string, error) {
	results, err := r.client.Keys(ctx, fmt.Sprintf("*%s*", mustInclude)).Result()
	if err != nil {
		zap.L().Error("GetKeys error", zap.Error(err), zap.String("must_include", mustInclude))
		return nil, err
	}

	if len(results) == 0 {
		zap.L().Error("GetKeys error", zap.Error(redis.Nil), zap.String("must_include", mustInclude))
		return nil, redis.Nil
	}

	return results, nil
}

func (r *rdm) GetIncludingKey(ctx context.Context, mustInclude string) (string, error) {
	results, err := r.client.Keys(ctx, fmt.Sprintf("*%s*", mustInclude)).Result()
	if err != nil {
		zap.L().Error("GetIncludingKey error", zap.Error(err), zap.String("must_include", mustInclude))
		return "", err
	}

	if len(results) == 0 {
		zap.L().Error("GetIncludingKey error", zap.Error(redis.Nil), zap.String("must_include", mustInclude))
		return "", redis.Nil
	}

	result := "["

	for i, k := range results {
		s, err := r.client.Get(ctx, k).Result()
		if err != nil {
			zap.L().Error("GetIncludingKey error", zap.Error(err), zap.String("must_include", mustInclude))
			return "", err
		}

		if i == 0 {
			result += s

			continue
		}

		result += "," + s
	}

	return result + "]", nil
}

func (r *rdm) Close() error {
	return r.client.Close()
}
