package storage

import (
	"errors"
	"fmt"
	"log"
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
func (s *Storage) ReadGauge(name string) (value float64, err error) {
	s.Lock()
	value, ok := s.metrics[name]
	s.Unlock()

	if !ok {
		err = &models.RequestError{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("metric not found"),
		}

		return
	}

	return
}

func (s *Storage) ReadCounter(name string) (value int64, err error) {
	s.Lock()
	f, ok := s.metrics[name]
	s.Unlock()

	if !ok {
		err = &models.RequestError{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("metric not found"),
		}

		return
	}

	value = int64(f)

	return
}

func (s *Storage) SetOldGauge(metric models.MetricGauge) (err error) {
	s.Lock()
	if s.metrics == nil {
		s.metrics = make(map[string]float64)
	}

	s.metrics[metric.Name] = metric.Value
	s.Unlock()

	return
}

func (s *Storage) SetOldCounter(metric models.MetricCounter) (err error) {
	s.Lock()
	if s.metrics == nil {
		s.metrics = make(map[string]float64)
	}

	value := s.metrics[metric.Name]
	s.metrics[metric.Name] = value + float64(metric.Value)
	s.Unlock()

	return
}

func (s *Storage) SetMetric(metric models.Metric) (err error) {
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

	s.metrics[metric.ID] = value
	s.Unlock()

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
