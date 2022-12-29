package config

import "time"

type Config struct {
	DB      DBConfig
	Storage StorageConfig
	Cache   CacheConfig
}

type DBConfig struct {
	Implementation string
}

type StorageConfig struct {
	Implementation string
}

type CacheConfig struct {
	TicketTTL time.Duration
	QueueTTL  time.Duration
	UserTTL   time.Duration
}

var config Config

func init() {
	config = Config{
		DB: DBConfig{
			Implementation: "SQL",
		},
		Storage: StorageConfig{
			Implementation: "FS",
		},
		Cache: CacheConfig{
			TicketTTL: 10 * time.Minute,
			QueueTTL:  60 * time.Minute,
			UserTTL:   10 * time.Minute,
		},
	}
}

func GetConfig() Config {
	return config
}
