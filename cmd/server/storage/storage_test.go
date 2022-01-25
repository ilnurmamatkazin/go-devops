package storage

import (
	"reflect"
	"testing"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg models.Config
	}
	tests := []struct {
		name        string
		args        args
		wantStorage *Storage
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStorage, err := New(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotStorage, tt.wantStorage) {
				t.Errorf("New() = %v, want %v", gotStorage, tt.wantStorage)
			}
		})
	}
}
