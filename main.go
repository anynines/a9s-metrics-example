package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	rand.Seed(time.Now().UnixNano())
	max := 10

	go func() {
		for {
			opsProcessed.Inc()
			randomNumber.Set(rand.Float64() * float64(max))
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "a9s_metrics_sample_ops_total",
		Help: "The total number of processed events",
	})
	randomNumber = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "a9s_metrics_sample_random_number",
		Help: "The total number of processed events",
	})
)

func main() {
	recordMetrics()

	port := "3000"
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "3000"
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on :%v\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
