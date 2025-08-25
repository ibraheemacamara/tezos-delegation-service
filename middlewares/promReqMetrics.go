package middlewares

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	log "github.com/sirupsen/logrus"
)

func PromReqMetrics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		go promReqsCounter(ctx)
		ctx.Next()
	}
}

func promReqsCounter(ctx *gin.Context) {
	reqCount(ctx)
	reqStatusCodeCount(ctx)
}

// Count number of requests by IP
func reqCount(ctx *gin.Context) {
	monitor := ginmetrics.GetMonitor()
	metricsName := "request_count_by_ip"
	ip := ctx.ClientIP()
	reqGauge := monitor.GetMetric(metricsName)
	if reqGauge.Name == "" {
		gaugeMetric := &ginmetrics.Metric{
			Type:        ginmetrics.Gauge,
			Name:        metricsName,
			Description: "Number of requests by IP address",
			Labels:      []string{"ip"},
		}
		err := monitor.AddMetric(gaugeMetric)
		if err != nil {
			log.Errorf("failed to add metrics %v: %v", metricsName, err)
		}
	}
	err := monitor.GetMetric(metricsName).Inc([]string{ip})
	if err != nil {
		log.Errorf("error while incrementing metric %v: %v", metricsName, err)
	}
}

// Count status code response by IP
func reqStatusCodeCount(ctx *gin.Context) {
	ip := ctx.ClientIP()
	statusCode := ctx.Writer.Status()
	monitor := ginmetrics.GetMonitor()
	metricsName := "request_status_code"
	metric := monitor.GetMetric(metricsName)
	if metric.Name == "" {
		gaugeMetric := &ginmetrics.Metric{
			Type:        ginmetrics.Gauge,
			Name:        metricsName,
			Description: "Response status code by id address",
			Labels:      []string{"statusCode", "ip"},
		}
		err := monitor.AddMetric(gaugeMetric)
		if err != nil {
			log.Errorf("failed to add metric %v. %v", metricsName, err)
		}
	}
	statusCodeStr := strconv.Itoa(statusCode)
	err := monitor.GetMetric(metricsName).Inc([]string{statusCodeStr, ip})
	if err != nil {
		log.Errorf("error while incrementing metric %v. %v", metricsName, err)
	}
}
