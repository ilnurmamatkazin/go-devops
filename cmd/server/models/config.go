package models

type Config struct {
	Address       string `env:"ADDRESS"`
	StoreInterval string `env:"STORE_INTERVAL"`
	StoreFile     string `env:"STORE_FILE"`
	Restore       bool   `env:"RESTORE"`
	Key           string `env:"KEY"`
	Database      string `env:"DATABASE_DSN"`
}
