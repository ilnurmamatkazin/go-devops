package memory

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// )

// func (mr *MemoryRepository) SaveToFile() (err error) {
// 	fmt.Println("mr.fileName", mr.fileName)
// 	mr.file, err = os.OpenFile(mr.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	defer mr.closeFile()

// 	mr.Lock()
// 	if mr.repository == nil {
// 		mr.repository = make(map[string]float64)
// 	}

// 	if err = json.NewEncoder(mr.file).Encode(&mr.repository); err != nil {
// 		mr.Unlock()

// 		fmt.Println(err.Error())
// 		return
// 	}
// 	mr.Unlock()

// 	return
// }

// func (mr *MemoryRepository) loadFromFile() (err error) {
// 	mr.file, err = os.OpenFile(mr.fileName, os.O_RDONLY|os.O_CREATE, 0777)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	defer mr.closeFile()

// 	fmt.Println(mr.file)
// 	fmt.Println(mr.repository)

// 	mr.Lock()

// 	if mr.repository == nil {
// 		mr.repository = make(map[string]float64)
// 	}

// 	if err = json.NewDecoder(mr.file).Decode(&mr.repository); err != nil {
// 		mr.Unlock()
// 		fmt.Println(err.Error())

// 		return
// 	}
// 	mr.Unlock()

// 	return
// }

// //Функция нужна если будем делать логирование
// func (mr *MemoryRepository) closeFile() (err error) {
// 	if err = mr.file.Close(); err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	return
// }
