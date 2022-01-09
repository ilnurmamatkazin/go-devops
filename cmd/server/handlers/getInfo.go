package handlers

import (
	"fmt"
	"net/http"
)

func getInfo(w http.ResponseWriter, r *http.Request) {
	mutexCounter.Lock()
	valueInt := storageCounter["PollCount"]
	mutexCounter.Unlock()

	mutexGauge.Lock()
	html := fmt.Sprintf(`
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
					<li>PollCount: %v</li>
				</ul>
			</body>
	</html>`,
		storageGauge["Alloc"],
		storageGauge["BuckHashSys"],
		storageGauge["Frees"],
		storageGauge["GCCPUFraction"],
		storageGauge["GCSys"],
		storageGauge["HeapAlloc"],
		storageGauge["HeapIdle"],
		storageGauge["HeapInuse"],
		storageGauge["HeapObjects"],
		storageGauge["HeapReleased"],
		storageGauge["HeapSys"],
		storageGauge["LastGC"],
		storageGauge["Lookups"],
		storageGauge["MCacheInuse"],
		storageGauge["MCacheSys"],
		storageGauge["MSpanInuse"],
		storageGauge["MSpanSys"],
		storageGauge["Mallocs"],
		storageGauge["NextGC"],
		storageGauge["NumForcedGC"],
		storageGauge["NumGC"],
		storageGauge["OtherSys"],
		storageGauge["PauseTotalNs"],
		storageGauge["StackInuse"],
		storageGauge["StackSys"],
		storageGauge["Sys"],
		storageGauge["RandomValue"],
		valueInt,
	)
	mutexGauge.Unlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(html))
}
