package main

import (
	"reflect"
	"testing"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func Test_parseConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantCfg models.Config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCfg, err := parseConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCfg, tt.wantCfg) {
				t.Errorf("parseConfig() = %v, want %v", gotCfg, tt.wantCfg)
			}
		})
	}
}
