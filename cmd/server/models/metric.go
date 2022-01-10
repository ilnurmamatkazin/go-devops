package models

type MetricCounter struct {
	Name  string
	Value int
}

type MetricGauge struct {
	Name  string
	Value float64
}
