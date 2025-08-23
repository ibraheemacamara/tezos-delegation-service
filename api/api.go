package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ibraheemacara/tezos-delegation-service/config"
	"github.com/ibraheemacara/tezos-delegation-service/db"
	"github.com/penglongli/gin-metrics/ginmetrics"
	log "github.com/sirupsen/logrus"
)

func StartServer(cfg config.Config, db db.DBInterface) {
	engine := gin.New()

	//metric
	metricRouter := gin.New()
	m := ginmetrics.GetMonitor()
	m.UseWithoutExposingEndpoint(engine)
	m.SetMetricPath("/metrics")
	m.Expose(metricRouter)
	go func() {
		log.Info(fmt.Sprintf("Metrics server started at url http://localhost:%v/metrics", cfg.Server.MetricsPort))

		_ = metricRouter.Run(fmt.Sprintf(":%v", cfg.Server.MetricsPort))
		log.Fatal("Metrics server stopped")
	}()

	ctrl := NewController(db)

	engine.GET("/delegations", ctrl.GetDelegations)
	engine.GET("/delegations/:year", ctrl.GetDelegations)

	engine.Run(fmt.Sprintf(":%v", cfg.Server.Port))
}
