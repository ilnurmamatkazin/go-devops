package file

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type FileRepository struct {
	file     *os.File
	FileName string
}

func (fr *FileRepository) SaveToFile(mutex *sync.Mutex, repository map[string]float64) (err error) {
	fr.file, err = os.OpenFile(fr.FileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer fr.closeFile()

	mutex.Lock()
	if repository == nil {
		repository = make(map[string]float64)
	}

	if err = json.NewEncoder(fr.file).Encode(&repository); err != nil {
		mutex.Unlock()

		log.Println(err.Error())
		return
	}
	mutex.Unlock()

	return
}

func (fr *FileRepository) LoadFromFile(mutex *sync.Mutex, repository map[string]float64) (err error) {
	fr.file, err = os.OpenFile(fr.FileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer fr.closeFile()

	mutex.Lock()

	if repository == nil {
		repository = make(map[string]float64)
	}

	if err = json.NewDecoder(fr.file).Decode(&repository); err != nil {
		mutex.Unlock()
		log.Println(err.Error())

		return
	}

	mutex.Unlock()

	return
}

//Функция нужна если будем делать логирование
func (fr *FileRepository) closeFile() (err error) {
	if err = fr.file.Close(); err != nil {
		log.Println(err.Error())
	}

	return
}
