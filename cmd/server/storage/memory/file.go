package memory

import (
	"encoding/json"
	"os"
)

func (mr *MemoryRepository) SaveToFile() (err error) {
	mr.file, err = os.OpenFile(mr.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return
	}

	defer mr.closeFile()

	mr.Lock()
	if err = json.NewEncoder(mr.file).Encode(&mr.repository); err != nil {
		return
	}
	mr.Unlock()

	return
}

func (mr *MemoryRepository) loadFromFile() (err error) {
	mr.file, err = os.OpenFile(mr.fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return
	}

	defer mr.closeFile()

	mr.Lock()
	if err = json.NewDecoder(mr.file).Decode(&mr.repository); err != nil {
		return
	}
	mr.Unlock()

	return
}

//Функция нужна если будем делать логирование
func (mr *MemoryRepository) closeFile() (err error) {
	err = mr.file.Close()
	return
}
