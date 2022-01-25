package file

import (
	"os"
	"sync"
	"testing"
)

func TestFileRepository_SaveToFile(t *testing.T) {
	type fields struct {
		file     *os.File
		FileName string
	}
	type args struct {
		mutex      *sync.Mutex
		repository map[string]float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := &FileRepository{
				file:     tt.fields.file,
				FileName: tt.fields.FileName,
			}
			if err := fr.SaveToFile(tt.args.mutex, tt.args.repository); (err != nil) != tt.wantErr {
				t.Errorf("FileRepository.SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
