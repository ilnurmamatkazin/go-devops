package pg

import (
	"testing"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	tests := []struct {
		name     string
		database string
		wantErr  bool
	}{
		{
			name:     "pisitive test",
			database: "postgres://postgres:12345@localhost:5434/postgres?sslmode=disable",
			wantErr:  false,
		},
		{
			name:     "negative test",
			database: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := models.Config{Database: tt.database}

			db, err := NewRepository(&cfg)
			if err != nil {
				db = nil
			}

			if db != nil {
				err = db.Ping()
				assert.Nil(t, err)
			}
		})
	}
}

func TestRepository_Init(t *testing.T) {
	type fields struct {
		conn *pgx.Conn
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				conn: tt.fields.conn,
			}

			err := r.Init()
			assert.Equal(t, (err != nil) == tt.wantErr, tt.wantErr)

		})
	}
}
