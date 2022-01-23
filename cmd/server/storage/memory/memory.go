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
	if err == nil {
		log.Fatalf("Ошибка создания тикера")
	}

	// strDurationStoreInterval := cfg.StoreInterval[len(cfg.StoreInterval)-1:]
	// strStoreInterval := cfg.StoreInterval[0 : len(cfg.StoreInterval)-1]

	// storeInterval, _ := strconv.Atoi(strStoreInterval)

	// var durationStoreInterval time.Duration

	// switch strDurationStoreInterval {
	// case "s":
	// 	durationStoreInterval = time.Second
	// case "m":
	// 	durationStoreInterval = time.Minute
	// case "h":
	// 	durationStoreInterval = time.Hour
	// default:
	// 	log.Println("Неверный временной интервал")
	// 	return memoryRepository
	// }

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
	mr.Lock()

	html = fmt.Sprintf(`
	<html>
		<head>
		<title></title>
		</head>
		<body>
			<ul>
				<li>Alloc: %f</li>
				<li>BuckHashSys: %f</li>
				<li>Frees: %f</li>
				<li>GCCPUFraction: %f</li>
				<li>GCSys: %f</li>
				<li>HeapAlloc: %f</li>
				<li>HeapIdle: %f</li>
				<li>HeapInuse: %f</li>
				<li>HeapObjects: %f</li>
				<li>HeapReleased: %f</li>
				<li>HeapSys: %f</li>
				<li>LastGC: %f</li>
				<li>Lookups: %f</li>
				<li>MCacheInuse: %f</li>
				<li>MCacheSys: %f</li>
				<li>MSpanInuse: %f</li>
				<li>MSpanSys: %f</li>
				<li>Mallocs: %f</li>
				<li>NextGC: %f</li>
				<li>NumForcedGC: %f</li>
				<li>NumGC: %f</li>
				<li>OtherSys: %f</li>
				<li>PauseTotalNs: %f</li>
				<li>TotalAlloc: %f</li>
				<li>StackInuse: %f</li>
				<li>StackSys: %f</li>
				<li>Sys: %f</li>
				<li>RandomValue: %f</li>
				<li>PollCount: %d</li>
			</ul>
		</body>
	</html>`,
		mr.repository["Alloc"],
		mr.repository["BuckHashSys"],
		mr.repository["Frees"],
		mr.repository["GCCPUFraction"],
		mr.repository["GCSys"],
		mr.repository["HeapAlloc"],
		mr.repository["HeapIdle"],
		mr.repository["HeapInuse"],
		mr.repository["HeapObjects"],
		mr.repository["HeapReleased"],
		mr.repository["HeapSys"],
		mr.repository["LastGC"],
		mr.repository["Lookups"],
		mr.repository["MCacheInuse"],
		mr.repository["MCacheSys"],
		mr.repository["MSpanInuse"],
		mr.repository["MSpanSys"],
		mr.repository["Mallocs"],
		mr.repository["NextGC"],
		mr.repository["NumForcedGC"],
		mr.repository["NumGC"],
		mr.repository["OtherSys"],
		mr.repository["PauseTotalNs"],
		mr.repository["TotalAlloc"],
		mr.repository["StackInuse"],
		mr.repository["StackSys"],
		mr.repository["Sys"],
		mr.repository["RandomValue"],
		int64(mr.repository["PollCount"]),
	)

	mr.Unlock()

	return
}
