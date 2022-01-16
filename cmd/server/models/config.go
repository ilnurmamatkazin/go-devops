package models

type Config struct {
	Address       string `env:"ADDRESS"`
	StoreInterval int    `env:"STORE_INTERVAL"`
	StoreFile     string `env:"STORE_FILE"`
	Restore       bool   `env:"RESTORE"`
}
