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

// InitRedis inicializa el cliente de Redis
func InitRedis(address, password string, db int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // Deja vacío si no se utiliza autenticación
		DB:       db,       // Número de la base de datos
	})

	// Opcional: Prueba la conexión a Redis para asegurar que está funcionando
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
}

// SetURL almacena una URL en el caché de Redis de forma reactiva
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

// GetURL recupera una URL o valor del caché de Redis de forma reactiva
func GetURL(key string) rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		result, err := rdb.Get(ctx, key).Result()
		if errors.Is(err, redis.Nil) {
			ch <- rxgo.Of("") // No encontrado en caché
		} else if err != nil {
			ch <- rxgo.Error(err)
		} else {
			ch <- rxgo.Of(result)
		}
	}})
}

// DeleteURL elimina una URL del caché de Redis de forma reactiva
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

// IncrementURLCounter incrementa el contador de accesos para una URL corta específica en Redis
func IncrementURLCounter(key string) error {
	_, err := rdb.Incr(ctx, key).Result()
	return err
}

// SetLastAccess establece la marca de tiempo del último acceso para una URL corta específica en Redis
func SetLastAccess(key string, timestamp string) error {
	return rdb.Set(ctx, key, timestamp, 0).Err()
}

// GetLastAccess recupera la marca de tiempo del último acceso para una URL corta específica de Redis
func GetLastAccess(key string) (string, error) {
	lastAccess, err := rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// Si la clave no existe, devuelve una cadena vacía
		return "", nil
	}
	return lastAccess, err
}
