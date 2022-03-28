package models

// Список констант, используемых в системе.
const (
	Address       = "localhost:8080"
	StoreInterval = "10s"
	StoreFile     = "/tmp/devops-metrics-db.json"
	Restore       = true
	Key           = ""
	// Database      = "postgres://postgres:1qaz@WSX@localhost:5432/postgres?sslmode=disable"
	Database        = "postgres://postgres:12345@localhost:5434/postgres?sslmode=disable"
	DatabaseTimeout = 30
	PrivateKey      = "../keys/private.pem"
	JSONConfig      = "./config.json"
)
