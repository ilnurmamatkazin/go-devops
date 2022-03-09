package pg

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// PgxIface интерфейс, необходим для абстагирования работы с драйвером СУБД.
type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
	Prepare(context.Context, string, string) (*pgconn.StatementDescription, error)
	Close(context.Context) error
}

// Repository структура для работы с методами взаимодействия с базой данных.
type Repository struct {
	conn PgxIface
}

// NewRepository конструктор для создания экземпляра структуры Repository.
func NewRepository(cfg *models.Config) (repository *Repository, err error) {
	repository = &Repository{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(models.DatabaseTimeout)*time.Second)
	defer cancel()

	if repository.conn, err = pgx.Connect(ctx, cfg.Database); err != nil {
		return
	}

	if err = repository.Init(); err != nil {
		return
	}

	return
}

// Init функция создания таблицы metrics в базе данных, в случае ее отсутсвия.
func (r *Repository) Init() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(models.DatabaseTimeout)*time.Second)
	defer cancel()

	query := `
	CREATE TABLE IF NOT EXISTS public.metrics
	(
		id text COLLATE pg_catalog."default" NOT NULL,
		type text COLLATE pg_catalog."default" NOT NULL,
		delta bigint,
		value double precision,
		hash text COLLATE pg_catalog."default",
		CONSTRAINT metrics_pkey PRIMARY KEY (id)
	)
	`

	if _, err = r.conn.Exec(ctx, query); err != nil {
		return
	}

	return
}

// Close функция закрития соединения с базой данных.
func (r *Repository) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(models.DatabaseTimeout)*time.Second)
	defer cancel()

	r.conn.Close(ctx)
}

// Ping функция проверки соединения с базой данных.
func (r *Repository) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(models.DatabaseTimeout)*time.Second)
	defer cancel()

	if r.conn != nil {
		return r.conn.Ping(ctx)
	} else {
		return errors.New("соединение с бд отсутствует")
	}

}

// Load функция проверки соединения с базой данных.
func (r *Repository) Load(mutex *sync.RWMutex, metrics map[string]models.Metric) (err error) {
	var (
		id, metricType string
		delta          sql.NullInt64
		value          sql.NullFloat64
		hash           sql.NullString
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(models.DatabaseTimeout)*time.Second)
	defer cancel()

	query := "select id, type, delta, value, hash from public.metrics"

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return
	}

	mutex.Lock()
	for rows.Next() {
		err = rows.Scan(&id, &metricType, &delta, &value, &hash)
		if err != nil {
			return
		}

		metric := models.Metric{ID: id, MetricType: metricType}

		if delta.Valid {
			i := delta.Int64
			metric.Delta = &i
		}

		if value.Valid {
			f := value.Float64
			metric.Value = &f
		}

		if hash.Valid {
			s := hash.String
			metric.Hash = &s
		}

		metrics[id] = metric
	}
	mutex.Unlock()

	return rows.Err()

}

// Save функция сохранения списка метрик из map в базе данных.
func (r *Repository) Save(mutex *sync.RWMutex, metrics map[string]models.Metric) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(models.DatabaseTimeout)*time.Second)
	defer cancel()

	query := `INSERT INTO public.metrics (id, type, delta, value, hash) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE SET delta=$3, value=$4, hash=$5`

	mutex.Lock()
	for key, value := range metrics {
		if _, err = r.conn.Exec(ctx, query, key, value.MetricType, value.Delta, value.Value, value.Hash); err != nil {
			return
		}
	}
	mutex.Unlock()

	r.Load(mutex, metrics)

	return

}

// SaveArray функция сохранения массива метрик в базе данных.
func (r *Repository) SaveArray(metrics []models.Metric) (err error) {
	if r.conn == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(models.DatabaseTimeout)*time.Second)
	defer cancel()

	query := `
	INSERT INTO public.metrics (id, type, delta, value, hash)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (id)
	DO UPDATE SET
	delta=$3,
	value=$4,
	hash=$5
	`

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return
	}

	defer tx.Rollback(ctx)

	for _, metric := range metrics {
		if _, err = r.conn.Exec(ctx, query, metric.ID, metric.MetricType, metric.Delta, metric.Value, metric.Hash); err != nil {
			return
		}
	}

	return tx.Commit(ctx)
}

// SaveCurentMetric функция сохранения метрики в базе данных.
func (r *Repository) SaveCurentMetric(metric models.Metric) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(models.DatabaseTimeout)*time.Second)
	defer cancel()

	query := `
	INSERT INTO public.metrics (id, type, delta, value, hash)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (id)
	DO UPDATE SET
	delta=$3,
	value=$4,
	hash=$5
	`

	if _, err = r.conn.Exec(ctx, query, metric.ID, metric.MetricType, metric.Delta, metric.Value, metric.Hash); err != nil {
		return
	}

	return

}
