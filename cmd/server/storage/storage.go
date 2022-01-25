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

type Storage struct {
	repository map[string]float64
	sync.Mutex
	isSyncMode bool
	Metric
}

func New(cfg models.Config) *Storage {
	mr := &Storage{
		repository: make(map[string]float64),
	}

	interval, _, err := utils.GetDataForTicker(cfg.StoreInterval)
	if err != nil {
		log.Fatalf("Ошибка получения параметров тикера: %s", err.Error())
	}

	if interval == 0 {
		mr.isSyncMode = true
	}

	_ = mr.initRepository(cfg)

	return mr
}

func (mr *Storage) initRepository(cfg models.Config) (err error) {
	if cfg.Database != "" {

	} else {
		if cfg.StoreFile == "" {
			return nil
		}

		fileRepository := &file.FileRepository{
			FileName: cfg.StoreFile,
		}

		if cfg.Restore {
			if err := fileRepository.LoadFromFile(&mr.Mutex, mr.repository); err != nil {
				log.Println(err.Error())
			}

			// defer fileRepository.CloseFile()
		}

		interval, duration, err := utils.GetDataForTicker(cfg.StoreInterval)
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

func (mr *Storage) ReadGauge(name string) (value float64, err error) {
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

func (mr *Storage) ReadCounter(name string) (value int64, err error) {
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

func (mr *Storage) SetGauge(metric models.MetricGauge) (err error) {
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

func (mr *Storage) SetCounter(metric models.MetricCounter) (err error) {
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

func (mr *Storage) Info() (html string) {
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

func (mr *Storage) Save() (err error) {
	return
}

func (mr *Storage) Load() (err error) {
	return
}
