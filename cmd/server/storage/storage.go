package storage

import "github.com/ilnurmamatkazin/go-devops/cmd/server/models"

type Metric interface {
	ReadGauge(name string) (float64, error)
	CreateGauge(metric models.MetricGauge) error
	ReadCounter(name string) (int, error)
	CreateCounter(metric models.MetricCounter) error
	Info() string
}

// type Storage struct {
// 	Metric
// }

// func NewStorage() *Repository {
// 	return &Repository{
// 		Authorization: NewAuthPostgres(db),
// 		TodoList:      NewTodoListPostgres(db),
// 		TodoItem:      NewTodoItemPostgres(db),
// 	}
// }
