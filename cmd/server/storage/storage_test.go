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
		name string
		args args
		want *Storage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
