package middleware

import (
	"context"
	"time"
	"userfc/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		timeOutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		ctx := context.WithValue(timeOutCtx, "request_id", requestID)
		c.Request = c.Request.WithContext(ctx)
		startTime := time.Now()
		c.Next()
		latency := time.Since(startTime)
		requestLog := logger.Fields{
			"request_field": requestID,
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"status":        c.Writer.Status(),
			"latency":       latency,
		}

		if c.Writer.Status() == 200 || c.Writer.Status() == 201 {
			logger.Logger.WithFields(requestLog).Info("request succes")

		} else {
			logger.Logger.WithFields(requestLog).Info("request error")
		}
	}
}
