package pg

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
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
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Conn: tt.fields.conn,
			}

			err := r.Init()
			assert.Equal(t, (err != nil) == tt.wantErr, tt.wantErr)

		})
	}
}

func TestRepository_Load(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	r := Repository{Conn: mock}

	type args struct {
		mutex   *sync.RWMutex
		metrics map[string]models.Metric
	}
	tests := []struct {
		name    string
		m       func()
		args    args
		wantErr bool
	}{
		{
			name: "Positive",
			m: func() {
				rows := mock.NewRows([]string{"id", "type", "delta", "value", "hash"}).
					AddRow("Alloc", "gauge", nil, 123.4, nil)
				mock.ExpectQuery("select id, type, delta, value, hash from public.metrics").WillReturnRows(rows)
			},
			args: args{
				metrics: make(map[string]models.Metric, 1),
				mutex:   new(sync.RWMutex),
			},
			wantErr: false,
		},
		{
			name: "Negative",
			m: func() {
				rows := mock.NewRows([]string{"id", "type", "delta", "value", "hash"}).
					AddRow("Alloc", "gauge", nil, 123.4, nil)
				mock.ExpectQuery("select id from public.metrics").WillReturnRows(rows)
			},
			args: args{
				metrics: make(map[string]models.Metric, 1),
				mutex:   new(sync.RWMutex),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m()

			err := r.Load(tt.args.mutex, tt.args.metrics)

			assert.Equal(t, (err != nil), tt.wantErr)

		})
	}
}

func TestRepository_Save(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	r := Repository{Conn: mock}

	type args struct {
		mutex   *sync.RWMutex
		metrics map[string]models.Metric
		metric  string
	}

	var metric models.Metric

	tests := []struct {
		name    string
		m       func()
		args    args
		wantErr bool
	}{
		{
			name: "Positive",
			m: func() {
				query := "^INSERT INTO public.metrics"

				mock.ExpectExec(query).WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			args: args{
				metrics: make(map[string]models.Metric, 1),
				mutex:   new(sync.RWMutex),
				metric:  `{"id": "PollCount", "type": "counter", "delta": 12345}`,
			},
			wantErr: false,
		},
		{
			name: "Negative",
			m: func() {
				query := `
					INSERT INTO public.metrics (id, type, delta, value, hash)
				`
				mock.ExpectExec(query)
			},
			args: args{
				metrics: make(map[string]models.Metric, 1),
				mutex:   new(sync.RWMutex),
				metric:  `{"id": "PollCount", "type": "counter", "delta": 12345}`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m()

			_ = json.Unmarshal([]byte(tt.args.metric), &metric)

			tt.args.mutex.Lock()
			tt.args.metrics[metric.ID] = metric
			tt.args.mutex.Unlock()

			err := r.Save(tt.args.mutex, tt.args.metrics)

			assert.Equal(t, (err != nil), tt.wantErr)
		})
	}
}
