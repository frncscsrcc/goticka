package config

type Config struct {
	DB      DBConfig
	Storage StorageConfig
}

type DBConfig struct {
	Implementation string
}

type StorageConfig struct {
	Implementation string
}

func GetConfig() Config {
	return Config{
		DB: DBConfig{
			Implementation: "SQL",
		},
		Storage: StorageConfig{
			Implementation: "FS",
		},
	}
}
