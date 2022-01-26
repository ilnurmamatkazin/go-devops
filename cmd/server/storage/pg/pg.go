package pg

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/jackc/pgx/v4"
)

type Repository struct {
	conn *pgx.Conn
}

func New(cfg models.Config) (repository *Repository, err error) {
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

func (r *Repository) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(models.DatabaseTimeout)*time.Second)
	defer cancel()

	r.conn.Close(ctx)
}

func (r *Repository) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(models.DatabaseTimeout)*time.Second)
	defer cancel()

	return r.conn.Ping(ctx)
}

func (r *Repository) Load(mutex *sync.Mutex, metrics map[string]float64) (err error) {
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

		if metricType == "gauge" {
			metrics[id] = value.Float64
		} else {
			metrics[id] = float64(delta.Int64)
		}

	}
	mutex.Unlock()

	return rows.Err()

}

func (r *Repository) Save(mutex *sync.Mutex, metrics map[string]float64) (err error) {
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
	for key, value := range metrics {
		// fmt.Println(key, "gauge", value, value, "qqq")

		if _, err = r.conn.Exec(ctx, query, key, "gauge", value, value, ""); err != nil {
			return
		}
	}

	return

}
