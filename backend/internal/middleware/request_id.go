package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "request_id"

func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader("X-Request-ID")

		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx.Set(RequestIDKey, requestID)
		ctx.Header("X-Request-ID", requestID)

		newCtx := context.WithValue(ctx.Request.Context(), RequestIDKey, requestID)
		ctx.Request = ctx.Request.WithContext(newCtx)

		ctx.Next()
	}
}
