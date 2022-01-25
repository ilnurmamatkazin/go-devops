package pg

import (
	"context"
	"time"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/jackc/pgx/v4"
)

type PGRepository struct {
	db *pgx.Conn
}

func NewPGRepository(cfg models.Config) (pgRepository *PGRepository, err error) {
	pgRepository = &PGRepository{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(models.DatabaseTimeout)*time.Second)
	defer cancel()

	if pgRepository.db, err = pgx.Connect(ctx, cfg.Database); err != nil {
		return
	}

	return
}

func (r *PGRepository) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(models.DatabaseTimeout)*time.Second)
	defer cancel()

	r.db.Close(ctx)
}
