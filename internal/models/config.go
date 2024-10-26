package models

type Config struct {
	Port            string
	MongoURI        string
	MongoDBName     string
	MongoCollection string
	RedisAddress    string
	RedisPassword   string
	RedisDB         int
}
