package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"service-order-gateway/internal/domain/dto"
	"time"
)

func LogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		timie := time.Now()

		ctx.Next()

		logMiddleware := dto.RequestLog{
			AcessTime: timie,
			Latency:   time.Since(timie),
			ClientIP:  ctx.ClientIP(),
			Method:    ctx.Request.Method,
			Code:      ctx.Writer.Status(),
			Path:      ctx.Request.URL.Path,
			UserAgent: ctx.Request.UserAgent(),
		}

		switch {
		case ctx.Writer.Status() >= 500:
			log.Fatal().Any("Internal Server Error", logMiddleware)
		case ctx.Writer.Status() >= 400:
			log.Warn().Any("Client Error", logMiddleware).Msg("")
		case ctx.Writer.Status() >= 300:
			log.Warn().Any("redirection", logMiddleware).Msg("")
		case ctx.Writer.Status() >= 200:
			log.Info().Any("SUCCESS", logMiddleware).Msg("")
		case ctx.Writer.Status() >= 200:
			log.Info().Any("INFO", logMiddleware).Msg("")
		}

	}
}
