package storage

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage/file"
	"github.com/ilnurmamatkazin/go-devops/internal/utils"
)

type Storage1 struct {
	repository map[string]float64
	sync.Mutex
	isSyncMode bool
	cfg        models.Config
}

func New(cfg models.Config) *Storage1 {
	mr := &Storage1{
		repository: make(map[string]float64),
		cfg:        cfg,
	}

	interval, _, err := utils.GetDataForTicker(cfg.StoreInterval)
	if err != nil {
		log.Fatalf("Ошибка получения параметров тикера: %s", err.Error())
	}

	if interval == 0 {
		mr.isSyncMode = true
	}

	_ = mr.initRepository()

	return mr
}

func (mr *Storage1) initRepository() (err error) {
	if mr.cfg.Database != "" {

	} else {
		if mr.cfg.StoreFile == "" {
			return nil
		}

		fmt.Println("@@1@@", mr.cfg)

		fileRepository := &file.FileRepository{
			FileName: mr.cfg.StoreFile,
		}

		if mr.cfg.Restore {
			if err := fileRepository.LoadFromFile(&mr.Mutex, mr.repository); err != nil {
				log.Println(err.Error())
			}

			// defer fileRepository.CloseFile()
		}

		interval, duration, err := utils.GetDataForTicker(mr.cfg.StoreInterval)
		if err != nil {
			log.Fatalf("Ошибка получения параметров тикера: %s", err.Error())
		}

		if interval == 0 {
			return err
		}
		go func(fr *file.FileRepository) {
			var err error
			ticker := time.NewTicker(time.Duration(interval) * duration)

			for {
				<-ticker.C

				if err = fr.SaveToFile(&mr.Mutex, mr.repository); err != nil {
					log.Println(err.Error())
				}
			}

		}(fileRepository)

	}

	return
}

func (mr *Storage1) ReadGauge(name string) (value float64, err error) {
	mr.Lock()
	value = mr.repository[name]
	mr.Unlock()

	if value == 0 {
		err = &models.RequestError{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("metric not found"),
		}
	}

	return
}

func (mr *Storage1) ReadCounter(name string) (value int64, err error) {
	mr.Lock()
	value = int64(mr.repository[name])
	mr.Unlock()

	if value == 0 {
		err = &models.RequestError{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("metric not found"),
		}
	}

	return
}

func (mr *Storage1) SetGauge(metric models.MetricGauge) (err error) {
	if mr.repository == nil {
		mr.Lock()
		mr.repository = make(map[string]float64)
		mr.Unlock()
	}

	mr.Lock()
	mr.repository[metric.Name] = metric.Value
	mr.Unlock()

	if mr.isSyncMode {
		err = mr.Save()
	}

	return
}

func (mr *Storage1) SetCounter(metric models.MetricCounter) (err error) {
	if mr.repository == nil {
		mr.Lock()
		mr.repository = make(map[string]float64)
		mr.Unlock()
	}

	mr.Lock()
	value := mr.repository[metric.Name]
	mr.repository[metric.Name] = value + float64(metric.Value)
	mr.Unlock()

	if mr.isSyncMode {
		err = mr.Save()
	}

	return
}

func (mr *Storage1) Info() (html string) {
	ul := ""

	mr.Lock()
	for key, value := range mr.repository {
		ul = ul + fmt.Sprintf("<li>%s: %f</li>", key, value)
	}
	mr.Unlock()

	html = fmt.Sprintf(`
	<html>
		<head>
		<title></title>
		</head>
		<body>
			<ul>%s</ul>
		</body>
	</html>`, ul)

	return
}

func (mr *Storage1) Save() (err error) {
	if mr.cfg.Database != "" {

	} else {
		if mr.cfg.StoreFile == "" {
			return nil
		}

		fileRepository := &file.FileRepository{
			FileName: mr.cfg.StoreFile,
		}

		if err = fileRepository.SaveToFile(&mr.Mutex, mr.repository); err != nil {
			log.Println(err.Error())
		}

	}

	return
}

func (mr *Storage1) Load() (err error) {
	return
}
