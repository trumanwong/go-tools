package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/trumanwong/go-tools/helper"
	"net/http"
	"time"
)

type prometheusMiddleware struct {
	buckets         []float64
	notStatisticUri []string
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestSize     *prometheus.SummaryVec
	responseSize    *prometheus.SummaryVec
}

func NewPrometheusMiddleware(buckets []float64, notStatisticUri []string) Middleware {
	if buckets == nil {
		buckets = prometheus.ExponentialBuckets(0.1, 1.5, 5)
	}

	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Tracks the number of HTTP requests.",
		}, []string{"method", "host", "uri", "handler"},
	)
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Tracks the latencies for HTTP requests.",
			Buckets: buckets,
		},
		[]string{"method", "uri", "handler", "status"},
	)
	requestSize := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_size_bytes",
			Help: "Tracks the size of HTTP requests.",
		},
		[]string{"method", "uri", "handler"},
	)
	responseSize := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_response_size_bytes",
			Help: "Tracks the size of HTTP responses.",
		},
		[]string{"method", "uri", "handler", "status"},
	)
	prometheus.MustRegister(
		requestsTotal,
		requestDuration,
		requestSize,
		responseSize,
	)
	return &prometheusMiddleware{
		buckets:         buckets,
		notStatisticUri: notStatisticUri,
		requestsTotal:   requestsTotal,
		requestDuration: requestDuration,
		requestSize:     requestSize,
		responseSize:    responseSize,
	}
}

func (m *prometheusMiddleware) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		if m.notStatisticUri != nil && helper.InArray(ctx.Request.RequestURI, m.notStatisticUri) {
			return
		}

		status := fmt.Sprintf("%d", ctx.Writer.Status())
		m.requestsTotal.WithLabelValues(
			ctx.Request.Method,
			ctx.Request.Host,
			ctx.Request.RequestURI,
			ctx.HandlerName(),
		).Inc()
		m.requestDuration.WithLabelValues(
			ctx.Request.Method,
			ctx.Request.RequestURI,
			ctx.HandlerName(),
			status,
		).Observe(time.Since(start).Seconds())
		m.requestSize.WithLabelValues(
			ctx.Request.Method,
			ctx.Request.RequestURI,
			ctx.HandlerName(),
		).Observe(float64(computeApproximateRequestSize(ctx.Request)))
		m.responseSize.WithLabelValues(
			ctx.Request.Method,
			ctx.Request.RequestURI,
			ctx.HandlerName(),
			status,
		).Observe(float64(ctx.Writer.Size()))
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
