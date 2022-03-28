package models

// Config структура конфигурации сервера.
type Config struct {
	Address       string `json:"address" env:"ADDRESS"`
	StoreInterval string `json:"store_interval" env:"STORE_INTERVAL"`
	StoreFile     string `json:"store_file" env:"STORE_FILE"`
	Restore       bool   `json:"restore" env:"RESTORE"`
	Key           string `json:"hash_key" env:"KEY"`
	Database      string `json:"database_dsn" env:"DATABASE_DSN"`
	PrivateKey    string `json:"crypto_key" env:"CRYPTO_KEY"`
	Config        string `env:"CONFIG"`
}
