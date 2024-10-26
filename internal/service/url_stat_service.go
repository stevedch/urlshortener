package service

import (
	"context"
	"strconv"
	"time"

	"github.com/reactivex/rxgo/v2"

	"urlshortener/internal/cache"
)

type URLStatService struct{}

// NewURLStatService creates a new instance of URLStatService
func NewURLStatService() *URLStatService {
	return &URLStatService{}
}

func (s *URLStatService) GetURLStats(shortID string) rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		// Gets the access counter
		countObservable := cache.GetURL(shortID + ":access_count")
		countResult := <-countObservable.Observe()
		if countResult.E != nil {
			ch <- rxgo.Error(countResult.E)
			return
		}
		countStr := countResult.V.(string)

		// Converts the counter to an integer
		var count int
		if countStr == "" {
			count = 0
		} else {
			var err error
			count, err = strconv.Atoi(countStr)
			if err != nil {
				ch <- rxgo.Error(err)
				return
			}
		}

		// Gets the last access timestamp
		lastAccessObservable := cache.GetURL(shortID + ":last_access")
		lastAccessResult := <-lastAccessObservable.Observe()
		if lastAccessResult.E != nil {
			ch <- rxgo.Error(lastAccessResult.E)
			return
		}
		lastAccess := lastAccessResult.V.(string)
		if lastAccess == "" {
			lastAccess = "N/A"
		}

		stats := map[string]interface{}{
			"access_count": count,
			"last_access":  lastAccess,
		}
		ch <- rxgo.Of(stats)
	}})
}

func (s *URLStatService) RecordAccess(shortID string) rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		// Increments the access counter
		err := cache.IncrementURLCounter(shortID + ":access_count")
		if err != nil {
			ch <- rxgo.Error(err)
			return
		}

		// Updates the last access timestamp
		timestamp := time.Now().Format(time.RFC3339)
		err = cache.SetLastAccess(shortID+":last_access", timestamp)
		if err != nil {
			ch <- rxgo.Error(err)
			return
		}

		ch <- rxgo.Of(true)
	}})
}
