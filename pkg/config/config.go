package config

import "time"

type Config struct {
	DB      DBConfig
	Storage StorageConfig
	Cache   CacheConfig
	Secrets Secrets
}

type DBConfig struct {
	Implementation string
}

type StorageConfig struct {
	Implementation string
	BasePath       string
}

type CacheConfig struct {
	TicketTTL time.Duration
	QueueTTL  time.Duration
	UserTTL   time.Duration
	RoleTTL   time.Duration
}

type Secrets struct {
	JWTSecret string
	JWTTTL    time.Duration
}

var config Config

func init() {
	config = Config{
		DB: DBConfig{
			Implementation: "SQL",
		},
		Storage: StorageConfig{
			Implementation: "FS",
			BasePath:       "./shared/attachments/",
		},
		Cache: CacheConfig{
			TicketTTL: 10 * time.Minute,
			QueueTTL:  60 * time.Minute,
			UserTTL:   10 * time.Minute,
			RoleTTL:   60 * time.Minute,
		},
		Secrets: Secrets{
			JWTSecret: "___SECRET___",
			JWTTTL:    1 * time.Hour,
		},
	}
}

func GetConfig() Config {
	return config
}

func OverwriteConfig(newConfig Config) Config {
	config = newConfig
	return newConfig
}
