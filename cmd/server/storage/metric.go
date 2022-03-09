package storage

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage/pg"
	"github.com/ilnurmamatkazin/go-devops/internal/utils"
)

type StorageMetrick struct {
	isSyncMode bool
	db         *pg.Repository
	metrics    map[string]models.Metric
	cfg        *models.Config
	sync.RWMutex
}

func NewStorageMetric(cfg *models.Config, db *pg.Repository) *StorageMetrick {
	return &StorageMetrick{
		metrics: make(map[string]models.Metric),
		cfg:     cfg,
		db:      db,
	}
}

func (s *StorageMetrick) ConnectPG() (err error) {
	interval, duration, err := utils.GetDataForTicker(s.cfg.StoreInterval)
	if err != nil {
		log.Fatalf("Ошибка получения параметров тикера: %s", err.Error())
		return
	}

	if s.cfg.Restore {
		if err = s.db.Load(&s.RWMutex, s.metrics); err != nil {
			log.Println(err.Error())
			return
		}
	}

	if interval == 0 {
		s.isSyncMode = true
	} else {
		go func(s *StorageMetrick, i int, d time.Duration) {
			var err error
			ticker := time.NewTicker(time.Duration(i) * d)

			for {
				<-ticker.C

				if err = s.db.Save(&s.RWMutex, s.metrics); err != nil {
					log.Println(err.Error())
				}
			}

		}(s, interval, duration)

	}

	return
}

func (s *StorageMetrick) Ping() error {
	return s.db.Ping()
}

func (s *StorageMetrick) ReadMetric(metric *models.Metric) (err error) {
	s.RLock()
	m, ok := s.metrics[metric.ID]
	s.RUnlock()

	if !ok {
		err = &models.RequestError{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("metric not found"),
		}

		return
	}

	if m.Delta != nil {
		metric.Delta = m.Delta
	}

	if m.Value != nil {
		metric.Value = m.Value
	}

	return
}

func (s *StorageMetrick) SetOldMetric(metric models.Metric) {
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
		metric.Value = nil
	} else {
		metric.Delta = nil
	}

	s.metrics[metric.ID] = metric

	s.Unlock()
}

func (s *StorageMetrick) SetMetric(metric models.Metric) (err error) {
	s.SetOldMetric(metric)

	if strings.Contains(metric.ID, "PopulateCounter") {
		if err = s.db.SaveCurentMetric(metric); err != nil {
			log.Println(err.Error())
			return
		}
	}

	if s.isSyncMode {
		if err = s.db.Save(&s.RWMutex, s.metrics); err != nil {
			log.Println(err.Error())
			return
		}
	}

	return
}

func (s *StorageMetrick) Info() (html string) {
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

func (s *StorageMetrick) SetArrayMetrics(metrics []models.Metric) (err error) {
	if err = s.db.SaveArray(metrics); err != nil {
		log.Println(err.Error())
		return
	}

	for _, metric := range metrics {
		s.SetOldMetric(metric)
	}
	return
}

func (s *StorageMetrick) Save() (err error) {
	return s.db.Save(&s.RWMutex, s.metrics)
}
