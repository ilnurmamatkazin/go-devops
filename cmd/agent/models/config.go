package models

// Config структура конфигурации агента.
type Config struct {
	Address        string `env:"ADDRESS"`         // адрес принимающего сервера
	ReportInterval string `env:"REPORT_INTERVAL"` // период отправки метрик
	PollInterval   string `env:"POLL_INTERVAL"`   // период сбора метрик
	Key            string `env:"KEY"`             // ключ для формирования подписи
}
