package handlers

import (
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"net/http"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func (h *Handler) getMetric(w http.ResponseWriter, r *http.Request) {
	var (
		metric models.Metric
		err    error
	)

	// fmt.Println("&&&&increment11 getMetric&&&")

	if err = json.NewDecoder(r.Body).Decode(&metric); err != nil {
		// fmt.Println("&&&& increment11 getMetric err &&&", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// fmt.Println("&&&&increment11 getMetric&&&", metric)

	if err = h.service.GetMetric(&metric); err != nil {
		re, ok := err.(*models.RequestError)
		if ok {
			http.Error(w, re.Err.Error(), re.StatusCode)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	// var (
	// 	i int64   = -100
	// 	f float64 = -100
	// )

	// if metric.Delta != nil {
	// 	i = *metric.Delta
	// }

	// if metric.Value != nil {
	// 	f = *metric.Value
	// }

	// fmt.Println("&&&&increment11 getMetric&&&", metric, i, f)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) parseMetric(w http.ResponseWriter, r *http.Request) {
	var (
		metric models.Metric
		err    error
	)

	fmt.Println("&&&&increment11 parseMetric&&&")

	if err = json.NewDecoder(r.Body).Decode(&metric); err != nil {
		fmt.Println("&&&& increment11 parseMetric err &&&", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.service.SetMetric(metric); err != nil {
		re, ok := err.(*models.RequestError)
		if ok {
			http.Error(w, re.Err.Error(), re.StatusCode)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	var (
		i int64   = -100
		f float64 = -100
	)

	if metric.Delta != nil {
		i = *metric.Delta
	}

	if metric.Value != nil {
		f = *metric.Value
	}

	fmt.Println("&&&&increment11 parseMetric&&&", metric, i, f)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// func sendMetricJSONData(w http.ResponseWriter, metric models.Metric) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	if err := json.NewEncoder(w).Encode(metric); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

func (h *Handler) parseMetrics(w http.ResponseWriter, r *http.Request) {
	var (
		metrics []models.Metric
		err     error
	)

	if err = json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		fmt.Println("increment11 parseMetrics Decode err: ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("increment11 parseMetrics metrics: ", metrics)

	if err = h.service.SetArrayMetrics(metrics); err != nil {
		fmt.Println("increment11 parseMetrics SetArrayMetrics err: ", err.Error())

		re, ok := err.(*models.RequestError)
		if ok {
			http.Error(w, re.Err.Error(), re.StatusCode)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		fmt.Println("increment11 parseMetrics Encode err: ", err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("increment11 parseMetrics http.StatusOK: ", http.StatusOK)

}
