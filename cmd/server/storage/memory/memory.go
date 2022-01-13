package memory

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

type MemoryRepository struct {
	counter map[string]int
	gauge   map[string]float64
	sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		counter: make(map[string]int),
		gauge:   make(map[string]float64),
	}
}

func (mr *MemoryRepository) ReadGauge(name string) (value float64, err error) {
	mr.Lock()
	value = mr.gauge[name]
	mr.Unlock()

	if value == 0 {
		err = &models.RequestError{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("metric not found"),
		}
	}

	return
}

func (mr *MemoryRepository) ReadCounter(name string) (value int, err error) {
	mr.Lock()
	value = mr.counter[name]
	mr.Unlock()

	if value == 0 {
		err = &models.RequestError{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("metric not found"),
		}
	}

	return
}

func (mr *MemoryRepository) CreateGauge(metric models.MetricGauge) (err error) {
	if mr.gauge == nil {
		mr.Lock()
		mr.gauge = make(map[string]float64)
		mr.Unlock()
	}

	// if _, ok := mr.customers[c.GetID()]; ok {
	// 	return fmt.Errorf("customer already exists: %w", customer.ErrFailedToAddCustomer)
	// }
	mr.Lock()
	mr.gauge[metric.Name] = metric.Value
	mr.Unlock()

	return
}

func (mr *MemoryRepository) CreateCounter(metric models.MetricCounter) (err error) {
	if mr.counter == nil {
		mr.Lock()
		mr.counter = make(map[string]int)
		mr.Unlock()
	}

	mr.Lock()
	value := mr.counter[metric.Name]
	mr.counter[metric.Name] = value + metric.Value
	mr.Unlock()

	return
}

func (mr *MemoryRepository) Info() (html string) {
	mr.Lock()
	valueInt := mr.counter["PollCount"]
	mr.Unlock()

	mr.Mutex.Lock()
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
				<li>StackInuse: %f</li>
				<li>StackSys: %f</li>
				<li>Sys: %f</li>
				<li>RandomValue: %f</li>
				<li>PollCount: %d</li>
			</ul>
		</body>
	</html>`,
		mr.gauge["Alloc"],
		mr.gauge["BuckHashSys"],
		mr.gauge["Frees"],
		mr.gauge["GCCPUFraction"],
		mr.gauge["GCSys"],
		mr.gauge["HeapAlloc"],
		mr.gauge["HeapIdle"],
		mr.gauge["HeapInuse"],
		mr.gauge["HeapObjects"],
		mr.gauge["HeapReleased"],
		mr.gauge["HeapSys"],
		mr.gauge["LastGC"],
		mr.gauge["Lookups"],
		mr.gauge["MCacheInuse"],
		mr.gauge["MCacheSys"],
		mr.gauge["MSpanInuse"],
		mr.gauge["MSpanSys"],
		mr.gauge["Mallocs"],
		mr.gauge["NextGC"],
		mr.gauge["NumForcedGC"],
		mr.gauge["NumGC"],
		mr.gauge["OtherSys"],
		mr.gauge["PauseTotalNs"],
		mr.gauge["StackInuse"],
		mr.gauge["StackSys"],
		mr.gauge["Sys"],
		mr.gauge["RandomValue"],
		valueInt,
	)
	mr.Mutex.Unlock()

	return
}
