package cache

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/reactivex/rxgo/v2"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

// InitRedis initializes the Redis client and returns the instance
func InitRedis(address, password string, db int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // Leave empty if no authentication is used
		DB:       db,       // Database number
	})

	// Optional: Test the connection to Redis to ensure it's working
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
}

// SetURL stores a URL in the Redis cache reactively
func SetURL(shortID string, originalURL string) rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		err := rdb.Set(ctx, shortID, originalURL, 24*time.Hour).Err()
		if err != nil {
			ch <- rxgo.Error(err)
		} else {
			ch <- rxgo.Of(shortID)
		}
	}})
}

// GetURL retrieves a URL from the Redis cache reactively
func GetURL(shortID string) rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		result, err := rdb.Get(ctx, shortID).Result()
		if errors.Is(err, redis.Nil) {
			ch <- rxgo.Of("") // Not found in cache
		} else if err != nil {
			ch <- rxgo.Error(err)
		} else {
			ch <- rxgo.Of(result)
		}
	}})
}

// DeleteURL removes a URL from the Redis cache reactively
func DeleteURL(shortID string) rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		err := rdb.Del(ctx, shortID).Err()
		if err != nil {
			ch <- rxgo.Error(err)
		} else {
			ch <- rxgo.Of(shortID)
		}
	}})
}
