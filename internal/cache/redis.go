package cache

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	l "github.com/mathehluiz/plant-care-tracker/pkg/logger"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type ConnectionStorer interface {
	IsValidInterface() bool
	Close() error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Set(ctx context.Context, duration time.Duration, key string, value string) error
	GetKeys(ctx context.Context, mustInclude string) ([]string, error)
	GetIncludingKey(ctx context.Context, mustInclude string) (string, error)
}

var (
	strRedisErr = "Cannot start redis client"
	strParseErr = "Cannot parse data"
)

var ErrNil = redis.Nil

func Start(ctx context.Context) (ConnectionStorer, error) {
	host := os.Getenv("CACHE_HOST")
	strDbInstance := os.Getenv("CACHE_INSTANCE")

	l.Logger.Debug("connecting to redis", zap.String("host", host))

	if host == "" {
		l.Logger.Fatal(strRedisErr, zap.Error(errors.New("invalid host")))
	}

	dbInstance := 0
	if i, err := strconv.Atoi(strDbInstance); err == nil {
		dbInstance = i
	}

	poolSize := "10"
	parsedPoolSize, err := strconv.Atoi(poolSize)
	if err != nil {
		l.Logger.Fatal(strParseErr, zap.Error(err))
	}

	maxRetries := "3"
	parsedMaxRetries, err := strconv.Atoi(maxRetries)
	if err != nil {
		l.Logger.Fatal(strParseErr, zap.Error(err))
	}

	readTimeout := "20000"
	parsedreadTimeout, err := strconv.Atoi(readTimeout)
	if err != nil {
		l.Logger.Fatal(strParseErr, zap.Error(err))
	}

	writeTimeout := "20000"
	parsedWriteTimeout, err := strconv.Atoi(writeTimeout)
	if err != nil {
		l.Logger.Fatal(strParseErr, zap.Error(err))
	}

	conn := redis.NewClient(&redis.Options{
		Addr:         host,
		DB:           dbInstance,
		PoolSize:     parsedPoolSize,
		MaxRetries:   parsedMaxRetries,
		ReadTimeout:  time.Duration(parsedreadTimeout) * time.Second,
		WriteTimeout: time.Duration(parsedWriteTimeout) * time.Second,
	})

	ctx, cancel := context.WithTimeout(ctx, time.Duration(parsedreadTimeout)*time.Second)
	defer cancel()

	if err = conn.Ping(ctx).Err(); err != nil {
		l.Logger.Fatal(strRedisErr, zap.Error(err))
	}

	l.Logger.Debug("connected successfully to redis", zap.String("host", host))

	return &rdm{client: conn}, nil
}

func StartMock() (ConnectionStorer, error) {
	mClient, err := miniredis.Run()
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: mClient.Addr(),
	})

	return &rdm{client: redisClient}, nil
}
