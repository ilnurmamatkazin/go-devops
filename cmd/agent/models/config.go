package models

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval string `env:"REPORT_INTERVAL"`
	PollInterval   string `env:"POLL_INTERVAL"`
	Key            string `env:"KEY"`
}
