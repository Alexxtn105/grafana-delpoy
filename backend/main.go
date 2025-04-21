package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Кастомные метрики Prometheus
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "my_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "my_http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.1, 0.3, 0.5, 1.0, 2.5, 5.0},
		},
		[]string{"method", "path"},
	)

	httpErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "my_http_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"method", "path", "status"},
	)
)

func init() {
	// Регистрируем все кастомные метрики
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpErrorsTotal)
}

func main() {
	// Настройка логгера
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting server...")

	// Создаем Gin-роутер без стандартного middleware
	router := gin.New()

	// Добавляем middleware для логирования и восстановления после паник
	router.Use(gin.Logger(), gin.Recovery())

	// Регистрируем обработчики
	router.GET("/", rootHandler)
	router.GET("/status/:status_code", statusHandler)
	router.GET("/metrics", metricsHandler())

	// Запускаем сервер
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Обработчик метрик Prometheus
func metricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func rootHandler(c *gin.Context) {
	start := time.Now()

	c.String(http.StatusOK, "root")

	// Обновляем метрики
	httpRequestsTotal.WithLabelValues(
		c.Request.Method,
		c.Request.URL.Path,
		strconv.Itoa(http.StatusOK),
	).Inc()

	httpRequestDuration.WithLabelValues(
		c.Request.Method,
		c.Request.URL.Path,
	).Observe(time.Since(start).Seconds())
}

func statusHandler(c *gin.Context) {
	start := time.Now()

	// Парсим параметры
	statusCode, err := strconv.Atoi(c.Param("status_code"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status code"})
		return
	}

	secondsSleep, _ := strconv.Atoi(c.Query("seconds_sleep"))

	// Логируем запрос
	log.Printf("Processing request: status_code=%d, seconds_sleep=%d", statusCode, secondsSleep)

	// Имитируем задержку если нужно
	if secondsSleep > 0 {
		time.Sleep(time.Duration(secondsSleep) * time.Second)
	}

	// Возвращаем ошибку если статус код не 200
	if statusCode != http.StatusOK {
		log.Printf("Returning error status: %d", statusCode)
		httpErrorsTotal.WithLabelValues(
			c.Request.Method,
			c.Request.URL.Path,
			strconv.Itoa(statusCode),
		).Inc()

		c.AbortWithStatusJSON(statusCode, gin.H{"error": "an error occurred"})
		return
	}

	// Обновляем метрики
	httpRequestsTotal.WithLabelValues(
		c.Request.Method,
		c.Request.URL.Path,
		strconv.Itoa(statusCode),
	).Inc()

	httpRequestDuration.WithLabelValues(
		c.Request.Method,
		c.Request.URL.Path,
	).Observe(time.Since(start).Seconds())

	c.JSON(http.StatusOK, gin.H{"data": "Hello"})
}
