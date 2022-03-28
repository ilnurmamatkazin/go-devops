package models

// Config структура конфигурации агента.
type Config struct {
	Address        string `json:"address" env:"ADDRESS"`                 // адрес принимающего сервера
	ReportInterval string `json:"report_interval" env:"REPORT_INTERVAL"` // период отправки метрик
	PollInterval   string `json:"poll_interval" env:"POLL_INTERVAL"`     // период сбора метрик
	Key            string `json:"hash_key" env:"KEY"`                    // ключ для формирования подписи
	PublicKey      string `json:"crypto_key" env:"CRYPTO_KEY"`           // открытый ключ для шифрования
	Config         string `env:"CONFIG"`                                 // имя json файла с конфигурацией
}
