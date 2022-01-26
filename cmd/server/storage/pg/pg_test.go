package pg

import (
	"testing"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg models.Config
	}
	tests := []struct {
		name           string
		args           args
		wantRepository *Repository
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRepository, err := New(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, gotRepository, tt.wantRepository)
		})
	}
}
