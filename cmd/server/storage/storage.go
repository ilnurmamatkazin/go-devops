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
	db         *pg.Repository
	metrics    map[string]models.Metric
	sync.Mutex
}

func New(cfg models.Config) (storage *Storage, err error) {
	storage = &Storage{metrics: make(map[string]models.Metric)}

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

func (s *Storage) ReadMetric(metric *models.Metric) (err error) {
	s.Lock()
	m, ok := s.metrics[metric.ID]
	s.Unlock()

	if !ok {
		fmt.Println("*****ReadOldMetric********", metric.ID)

		err = &models.RequestError{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("metric not found"),
		}

		return
	}

	metric = &m

	return
}

func (s *Storage) SetOldMetric(metric models.Metric) {
	s.Lock()
	if s.metrics == nil {
		s.metrics = make(map[string]models.Metric)
	}

	var delta int64

	if metric.MetricType == "counter" {
		if s.metrics[metric.ID].Delta == nil {
			delta = *metric.Delta
		} else {
			delta = *s.metrics[metric.ID].Delta + *metric.Delta
		}

		metric.Delta = &delta
	}

	s.metrics[metric.ID] = metric

	// fmt.Println("######", s.metrics[metric.ID], value)

	// if s.metrics[metric.ID] != value {
	// 	s.metrics[metric.ID] = value
	// } else {
	// 	s.metrics[metric.ID] = rand.Float64()
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
		if value.MetricType == "counter" {
			ul = ul + fmt.Sprintf("<li>%s: %d</li>", key, *value.Delta)
		} else {
			ul = ul + fmt.Sprintf("<li>%s: %f</li>", key, *value.Value)
		}

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
