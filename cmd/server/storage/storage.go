package storage

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage/pg"
	"github.com/ilnurmamatkazin/go-devops/internal/utils"
)

type Storage struct {
	isSyncMode bool
	// cfg        models.Config
	db      *pg.Repository
	metrics map[string]float64
	sync.Mutex
}

func New(cfg models.Config) (storage *Storage, err error) {
	// storage = &Storage{cfg: cfg}
	storage = &Storage{metrics: make(map[string]float64)}

	if storage.db, err = pg.New(cfg); err != nil {
		return
	}

	interval, duration, err := utils.GetDataForTicker(cfg.StoreInterval)
	if err != nil {
		log.Fatalf("Ошибка получения параметров тикера: %s", err.Error())
		return
	}

	if cfg.Restore {
		if err = storage.db.Load(&storage.Mutex, storage.metrics); err != nil {
			log.Println(err.Error())
			return
		}
	}

	if interval == 0 {
		storage.isSyncMode = true
	} else {
		go func(s *Storage, i int, d time.Duration) {
			var err error
			ticker := time.NewTicker(time.Duration(i) * d)

			for {
				<-ticker.C

				if err = s.db.Save(&storage.Mutex, storage.metrics); err != nil {
					log.Println(err.Error())
				}
			}

		}(storage, interval, duration)

	}

	return
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) ReadMetric(name string) (value float64, err error) {
	s.Lock()
	value, ok := s.metrics[name]
	s.Unlock()

	if !ok {
		fmt.Println("*****ReadOldMetric********", name)

		err = &models.RequestError{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("metric not found"),
		}

		return
	}

	return
}

func (s *Storage) SetOldMetric(metric models.Metric) {
	var value float64

	s.Lock()
	if s.metrics == nil {
		s.metrics = make(map[string]float64)
	}

	if metric.MetricType == "counter" {
		value = s.metrics[metric.ID] + float64(*metric.Delta)
	} else {
		value = *metric.Value
	}

	fmt.Println("######", s.metrics[metric.ID], value)

	s.metrics[metric.ID] = rand.Float64() * float64(rand.Int())

	// if s.metrics[metric.ID] != value {
	// 	s.metrics[metric.ID] = value
	// } else {
	// 	s.metrics[metric.ID] = rand.Float64()*rand.Int()
	// }
	s.Unlock()
}

func (s *Storage) SetMetric(metric models.Metric) (err error) {
	s.SetOldMetric(metric)

	if s.isSyncMode {
		if err = s.db.Save(&s.Mutex, s.metrics); err != nil {
			log.Println(err.Error())
			return
		}
	}

	return
}

func (s *Storage) Info() (html string) {
	ul := ""

	s.Lock()
	for key, value := range s.metrics {
		ul = ul + fmt.Sprintf("<li>%s: %f</li>", key, value)
	}
	s.Unlock()

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

func (s *Storage) Save() (err error) {
	if err = s.db.Save(&s.Mutex, s.metrics); err != nil {
		log.Println(err.Error())
		return
	}

	return
}

func (s *Storage) Ping() error {
	return s.db.Ping()
}
