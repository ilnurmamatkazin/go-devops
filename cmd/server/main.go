package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

const (
	url  = "127.0.0.1"
	port = 8080
)

var (
	mutex   sync.Mutex
	storage map[string]interface{}
)

func ParseMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	if ct := strings.ToLower(r.Header.Get("Content-Type")); ct != "text/plain; charset=utf-8" {
		http.Error(w, "Content-Type is not text/plain", http.StatusBadRequest)
		return
	}

	fmt.Println()

	arr := strings.Split(r.URL.Path, "/")
	elementCount := len(arr)

	fmt.Println(r.URL.Path)
	fmt.Println(arr)
	fmt.Println(elementCount)

	if elementCount != 5 {
		switch elementCount {
		case 1, 2, 3, 4:
			http.Error(w, "Request is not /update/type/name/value", http.StatusNotFound)
		default:
			http.Error(w, "Request is not /update/type/name/value", http.StatusBadRequest)
		}

		return
	}

	valueMetric := arr[4]
	nameMetric := arr[3]
	typeMetric := arr[2]

	fmt.Println("!!!!!", valueMetric, nameMetric, typeMetric)

	switch typeMetric {
	case "gauge":
		f, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s value is not float", nameMetric), http.StatusBadRequest)
			return
		}

		mutex.Lock()

		storage[nameMetric] = f

		mutex.Unlock()

	case "counter":
		var counter []int64

		i, err := strconv.ParseInt(valueMetric, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s value is not integer", nameMetric), http.StatusBadRequest)
			return
		}

		mutex.Lock()

		if storage[nameMetric] == nil {
			counter = make([]int64, 0)
		}

		counter = append(counter, i)
		storage[nameMetric] = counter

		mutex.Unlock()

	default:
		http.Error(w, "Metric type is not gaude or counter", http.StatusNotImplemented)
		return

	}

	fmt.Println("storage ", storage)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	storage = make(map[string]interface{})

	http.HandleFunc("/update/", ParseMetric)
	// http.HandleFunc("/update/gauge/BuckHashSys/", BuckHashSys)
	// http.HandleFunc("/update/gauge/Frees/", Frees)
	// http.HandleFunc("/update/gauge/GCCPUFraction/", GCCPUFraction)
	// http.HandleFunc("/update/gauge/GCSys/", GCSys)
	// http.HandleFunc("/update/gauge/HeapAlloc/", HeapAlloc)
	// http.HandleFunc("/update/gauge/HeapIdle/", HeapIdle)
	// http.HandleFunc("/update/gauge/HeapInuse/", HeapInuse)
	// http.HandleFunc("/update/gauge/HeapObjects/", HeapObjects)
	// http.HandleFunc("/update/gauge/HeapReleased/", HeapReleased)
	// http.HandleFunc("/update/gauge/HeapSys/", HeapSys)
	// http.HandleFunc("/update/gauge/LastGC/", LastGC)
	// http.HandleFunc("/update/gauge/Lookups/", Lookups)
	// http.HandleFunc("/update/gauge/MCacheInuse/", MCacheInuse)
	// http.HandleFunc("/update/gauge/MCacheSys/", MCacheSys)
	// http.HandleFunc("/update/gauge/MSpanInuse/", MSpanInuse)
	// http.HandleFunc("/update/gauge/MSpanSys/", MSpanSys)
	// http.HandleFunc("/update/gauge/Mallocs/", Mallocs)
	// http.HandleFunc("/update/gauge/NextGC/", NextGC)
	// http.HandleFunc("/update/gauge/NumForcedGC/", NumForcedGC)
	// http.HandleFunc("/update/gauge/NumGC/", NumGC)
	// http.HandleFunc("/update/gauge/OtherSys/", OtherSys)
	// http.HandleFunc("/update/gauge/PauseTotalNs/", PauseTotalNs)
	// http.HandleFunc("/update/gauge/StackInuse/", StackInuse)
	// http.HandleFunc("/update/gauge/StackSys/", StackSys)
	// http.HandleFunc("/update/gauge/Sys/", Sys)
	// http.HandleFunc("/update/counter/PollCount/", PollCount)
	// http.HandleFunc("/update/gauge/RandomValue/", RandomValue)

	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	fmt.Println("Server started...")

	<-quit

	fmt.Println("Server shutdown")
}
