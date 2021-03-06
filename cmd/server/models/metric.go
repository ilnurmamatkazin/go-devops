package models

// Metric структура, описывающая метрику в системе.
type Metric struct {
	ID         string   `json:"id"`              // имя метрики
	MetricType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta      *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value      *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash       *string  `json:"hash,omitempty"`  // значение хеш-функции
}
