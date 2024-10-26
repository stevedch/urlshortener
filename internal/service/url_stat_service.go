// internal/service/url_stat_service.go
package service

import (
	"context"
	"strconv"
	"time"

	"github.com/reactivex/rxgo/v2"

	"urlshortener/internal/cache"
)

type URLStatService struct{}

// NewURLStatService crea una nueva instancia de URLStatService
func NewURLStatService() *URLStatService {
	return &URLStatService{}
}

func (s *URLStatService) GetURLStats(shortID string) rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		// Obtiene el contador de accesos
		countObservable := cache.GetURL(shortID + ":access_count")
		countResult := <-countObservable.Observe()
		if countResult.E != nil {
			ch <- rxgo.Error(countResult.E)
			return
		}
		countStr := countResult.V.(string)

		// Convierte el contador a entero
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

		// Obtiene la última marca de tiempo de acceso
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
		// Incrementa el contador de accesos
		err := cache.IncrementURLCounter(shortID + ":access_count")
		if err != nil {
			ch <- rxgo.Error(err)
			return
		}

		// Actualiza la marca de tiempo del último acceso
		timestamp := time.Now().Format(time.RFC3339)
		err = cache.SetLastAccess(shortID+":last_access", timestamp)
		if err != nil {
			ch <- rxgo.Error(err)
			return
		}

		ch <- rxgo.Of(true)
	}})
}
