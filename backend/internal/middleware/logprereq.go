package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"slices"

	"github.com/andibalo/flip-test/internal/constants"
	"github.com/andibalo/flip-test/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var skipLogPayloadPaths = []string{
	"/upload",
}

func shouldSkipLogPayload(uriPath string) bool {
	return slices.Contains(skipLogPayloadPaths, uriPath)
}

func LogPreReq(logger logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var clientID = ctx.Request.Header.Get(constants.XClientID)
		payload, _ := io.ReadAll(ctx.Request.Body)

		compactPayload := &bytes.Buffer{}
		err := json.Compact(compactPayload, payload)
		if err != nil {
			compactPayload = bytes.NewBuffer(payload)
		}

		ctx.Set("x-client-id", clientID)
		ctx.Set("path", ctx.Request.URL.Path)
		ctx.Set("method", ctx.Request.Method)

		if shouldSkipLogPayload(ctx.Request.URL.Path) {
			logger.InfoWithContext(ctx, "Interceptor Log")

			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(payload))

			ctx.Next()

			return
		}

		logger.InfoWithContext(ctx, "Interceptor Log",
			zap.Any("payload", compactPayload),
		)

		ctx.Set("payload", string(payload))

		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(payload))

		ctx.Next()
	}
}
