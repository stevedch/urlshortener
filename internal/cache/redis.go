package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/reactivex/rxgo/v2"
	"time"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

// InitRedis initializes the Redis client
func InitRedis(address, password string, db int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // Leave empty if no authentication is used
		DB:       db,       // Database number
	})

	// Optional: Test the Redis connection to ensure it's working
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

// GetURL retrieves a URL or value from the Redis cache reactively
func GetURL(key string) rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		result, err := rdb.Get(ctx, key).Result()
		if errors.Is(err, redis.Nil) {
			ch <- rxgo.Of("") // Not found in cache
		} else if err != nil {
			ch <- rxgo.Error(err)
		} else {
			ch <- rxgo.Of(result)
		}
	}})
}

// DeleteURL deletes a URL from the Redis cache reactively
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

// IncrementURLCounter increments the access counter for a specific short URL in Redis
func IncrementURLCounter(key string) error {
	_, err := rdb.Incr(ctx, key).Result()
	return err
}

// SetLastAccess sets the timestamp of the last access for a specific short URL in Redis
func SetLastAccess(key string, timestamp string) error {
	return rdb.Set(ctx, key, timestamp, 0).Err()
}

// GetLastAccess retrieves the last access timestamp for a specific short URL from Redis
func GetLastAccess(key string) (string, error) {
	lastAccess, err := rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// If the key does not exist, return an empty string
		return "", nil
	}
	return lastAccess, err
}
