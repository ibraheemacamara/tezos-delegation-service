package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func LoggerHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		end := time.Now()
		log.Infof("Request %s %s processed in %s", ctx.Request.Method, ctx.Request.URL.Path, end.Sub(start))
	}
}
