package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
)

type prometheusMiddleware struct {
	buckets  []float64
	registry prometheus.Registerer
}

func NewPrometheusMiddleware(buckets []float64, registry prometheus.Registerer) Middleware {
	if buckets == nil {
		buckets = prometheus.ExponentialBuckets(0.1, 1.5, 5)
	}
	return &prometheusMiddleware{
		buckets:  buckets,
		registry: registry,
	}
}

func (m *prometheusMiddleware) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reg := prometheus.WrapRegistererWith(prometheus.Labels{"handler": ctx.HandlerName()}, m.registry)

		requestsTotal := promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Tracks the number of HTTP requests.",
			}, []string{"method", "code"},
		)
		requestDuration := promauto.With(reg).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Tracks the latencies for HTTP requests.",
				Buckets: m.buckets,
			},
			[]string{"method", "code"},
		)
		requestSize := promauto.With(reg).NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "http_request_size_bytes",
				Help: "Tracks the size of HTTP requests.",
			},
			[]string{"method", "code"},
		)
		responseSize := promauto.With(reg).NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "http_response_size_bytes",
				Help: "Tracks the size of HTTP responses.",
			},
			[]string{"method", "code"},
		)
		ctx.Next()

		status := fmt.Sprintf("%d", ctx.Writer.Status())
		requestsTotal.WithLabelValues(
			ctx.Request.Method,
			ctx.Request.RequestURI,
			status,
		).Inc()
		requestDuration.WithLabelValues(
			ctx.Request.Method,
			ctx.Request.RequestURI,
			status,
		)
		requestSize.WithLabelValues(
			ctx.Request.Method,
			ctx.Request.RequestURI,
		)
		responseSize.WithLabelValues(
			ctx.Request.Method,
			ctx.Request.RequestURI,
			status,
		)
	}
}

func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s += len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
