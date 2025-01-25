package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

const RequestIDKey = "X-Request-ID"

type ErrTypeCastFailed struct {
	msg string
}

func (e *ErrTypeCastFailed) Error() string {
	return e.msg
}

func NewErrTypeCastFailed(msg string) *ErrTypeCastFailed {
	return &ErrTypeCastFailed{msg: msg}
}

type ErrFallbackUuidUsed struct {
	msg string
}

func (e *ErrFallbackUuidUsed) Error() string {
	return e.msg
}

func NewErrFallbackUuidUsed(msg string) *ErrFallbackUuidUsed {
	return &ErrFallbackUuidUsed{msg: msg}
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := otel.Tracer("backend/RequestIDMiddleware")
		ctx, span := tracer.Start(c.Request.Context(), "RequestIDMiddleware")
		c.Request = c.Request.WithContext(ctx)
		defer span.End()

		requestID := c.GetHeader(RequestIDKey)
		if requestID == "" {
			span.AddEvent("no request id in header")
			requestID = genUuid()
		}

		c.Set(RequestIDKey, requestID)
		c.Writer.Header().Set(RequestIDKey, requestID)

		span.SetAttributes(attribute.String("requestID", requestID))

		c.Next()
	}
}

// GetRequestIDFromCtx is a helper function that extracts reqId from req context.
// GetRequestIDFromCtx can only return errors of type ErrFallbackUuidUsed or ErrTypeCastFailed.
// If ErrFallbackUuidUsed is returned, it means execution can proceed, but you should inform someone of that fact.
// If ErrTypeCastFailed is returned, it means uuid present in context but failure happened casting it to string.
func GetRequestIDFromCtx(c *gin.Context) (string, error) {
	var requestID string

	uuidFromCtx, ok := c.Get(RequestIDKey)
	if !ok {
		requestID = fallbackUuidHandler(c)
		return requestID, NewErrFallbackUuidUsed(
			fmt.Sprintf("X-Request-Id not found; generated new UUID: %s", requestID),
		)
	}

	requestID, ok = uuidFromCtx.(string)
	if !ok {
		return "", NewErrTypeCastFailed(
			fmt.Sprintf("cannot convert X-Request-Id to string: %s", uuidFromCtx),
		)
	}

	return requestID, nil
}

func fallbackUuidHandler(c *gin.Context) string {
	requestID := genUuid()
	c.Set(RequestIDKey, requestID)
	c.Writer.Header().Set(RequestIDKey, requestID)
	return requestID
}

func genUuid() string {
	return uuid.New().String()
}
