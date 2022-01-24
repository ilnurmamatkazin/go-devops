package memory

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/internal/utils"
)

type MemoryRepository struct {
	repository map[string]float64
	sync.Mutex
	file       *os.File
	fileName   string
	isSyncMode bool
}

func NewMemoryRepository(cfg models.Config) *MemoryRepository {
	memoryRepository := &MemoryRepository{
		repository: make(map[string]float64),
	}

	if cfg.StoreFile == "" {
		return nil
	}

	memoryRepository.fileName = cfg.StoreFile

	if cfg.Restore {
		if err := memoryRepository.loadFromFile(); err != nil {
			log.Println(err.Error())
		}
	}

	interval, duration, err := utils.GetDataForTicker(cfg.StoreInterval)
	if err != nil {
		log.Fatalf("Ошибка создания тикера")
	}

	if interval == 0 {
		memoryRepository.isSyncMode = true
	} else {
		go func(mr *MemoryRepository) {
			var err error
			ticker := time.NewTicker(time.Duration(interval) * duration)

			for {
				<-ticker.C

				if err = mr.SaveToFile(); err != nil {
					log.Println(err.Error())
				}
			}

		}(memoryRepository)
	}

	return memoryRepository
}

func (mr *MemoryRepository) ReadGauge(name string) (value float64, err error) {
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

func (mr *MemoryRepository) ReadCounter(name string) (value int64, err error) {
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

func (mr *MemoryRepository) SetGauge(metric models.MetricGauge) (err error) {
	if mr.repository == nil {
		mr.Lock()
		mr.repository = make(map[string]float64)
		mr.Unlock()
	}

	mr.Lock()
	mr.repository[metric.Name] = metric.Value
	mr.Unlock()

	if mr.isSyncMode {
		err = mr.SaveToFile()
	}

	return
}

func (mr *MemoryRepository) SetCounter(metric models.MetricCounter) (err error) {
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
		err = mr.SaveToFile()
	}

	return
}

func (mr *MemoryRepository) Info() (html string) {
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
