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

	for rows.Next() {

		err := rows.Scan(&id, &description)
		if err != nil {
			return err
		}
	}

	return rows.Err()

}

func (r *Repository) Save(mutex *sync.Mutex, metrics map[string]float64) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(models.DatabaseTimeout)*time.Second)
	defer cancel()

	// Сделать сохранение

	return

}
