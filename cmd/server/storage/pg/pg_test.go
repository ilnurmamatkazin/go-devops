package pg

import (
	"testing"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func TestNewPGRepository(t *testing.T) {
	type args struct {
		cfg models.Config
	}
	tests := []struct {
		name             string
		args             args
		wantPgRepository *PGRepository
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPGRepository(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPGRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(gotPgRepository, tt.wantPgRepository) {
			// 	t.Errorf("NewPGRepository() = %v, want %v", gotPgRepository, tt.wantPgRepository)
			// }
		})
	}
}
