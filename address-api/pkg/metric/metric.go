package metric

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"log"
	"strconv"
)

type Metrics interface {
	IncreaseHits(status int, method, path string)
	ObserveResponseTime(status int, method, path string, observeTime float64)
}

type metrics struct {
	HitsTotal prometheus.Counter
	Hits      *prometheus.CounterVec
	Times     *prometheus.HistogramVec
}

func (metric *metrics) IncreaseHits(status int, method, path string) {
	metric.HitsTotal.Inc()
	metric.Hits.WithLabelValues(strconv.Itoa(status), method, path).Inc()
}

func (metric *metrics) ObserveResponseTime(status int, method, path string, observeTime float64) {
	metric.Times.WithLabelValues(strconv.Itoa(status), method, path).Observe(observeTime)
}

func CreateMetrics(address, name string) (Metrics, error) {
	var metric metrics
	metric.HitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + "_hits_total",
	})
	if err := prometheus.Register(metric.HitsTotal); err != nil {
		return nil, err
	}

	metric.Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name + "_hits",
	}, []string{"status", "method", "path"})
	if err := prometheus.Register(metric.Hits); err != nil {
		return nil, err
	}

	go func() {
		app := fiber.New()
		app.Get("/metrics", func(ctx *fiber.Ctx) error {
			fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())(ctx.Context())
			return nil
		})

		log.Printf("Metrics server is running on port: %s", address)
		if err := app.Listen(address); err != nil {
			log.Fatal(err)
		}
	}()

	return &metric, nil
}
