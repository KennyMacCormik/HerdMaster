package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	defaultMaxRunning = 100
	defaultMaxWait    = 100
	defaultRetryAfter = 1
)

type RateLimiter struct {
	lg      *slog.Logger
	limiter chan struct{}

	maxRunning, maxWait, retryAfter    int
	running, total, timedOut, rejected atomic.Int32
}

func NewRateLimiter(maxRunning, maxWait, retryAfter int, lg *slog.Logger) *RateLimiter {
	maxRunning, maxWait, retryAfter = normalizeParams(maxRunning, maxWait, retryAfter, lg)
	rm := &RateLimiter{lg: lg, limiter: make(chan struct{}, maxRunning),
		maxRunning: maxRunning, maxWait: maxWait, retryAfter: retryAfter}
	return rm
}

func (rm *RateLimiter) GetRateLimiter() gin.HandlerFunc {

	return func(c *gin.Context) {
		rm.total.Add(1)
		defer rm.total.Add(-1)
		const (
			traceName = "gin.middleware.GetRateLimiter"
			spanName  = "rate limiting middleware"
		)
		// init trace
		tracer := otel.Tracer(traceName)
		ctx, span := tracer.Start(c.Request.Context(), spanName)
		c.Request = c.Request.WithContext(ctx)
		defer span.End()
		defer func(span trace.Span) {
			span.SetAttributes(
				attribute.Int("running", int(rm.running.Load())),
				attribute.Int("total", int(rm.total.Load())),
				attribute.Int("timedOut", int(rm.timedOut.Load())),
				attribute.Int("rejected", int(rm.rejected.Load())),
			)
		}(span)
		// init logger
		uuid, err := GetRequestIDFromCtx(c)
		if err != nil && errors.Is(err, &ErrTypeCastFailed{}) {
			span.AddEvent(
				"failed to get request ID from context",
				trace.WithAttributes(attribute.String("error", err.Error())),
			)
			rm.lg.Error("failed to get request ID from context", "error", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		reqLg := LogReq(c, uuid, rm.lg, true)
		// Log the Trace ID
		reqLg.Info("request trace ID", "traceID", span.SpanContext().TraceID().String())

		if err != nil {
			span.AddEvent(
				"fallback uuid used",
				trace.WithAttributes(attribute.String("error", err.Error())),
			)
			span.SetAttributes(attribute.String("fallback requestID", uuid))
			reqLg.Warn("fallback uuid used", "error", err.Error())
		}
		// reject if too may goroutines
		if rm.total.Load() >= int32(rm.maxWait) {
			rm.rejected.Add(1)
			span.AddEvent(
				"too many total requests, rejecting request",
				trace.WithAttributes(attribute.String("error", err.Error())),
			)
			reqLg.Error("too many total requests, rejecting request",
				"total", rm.total.Load(),
				"maxWait", rm.maxWait,
			)
			c.Header("Retry-After", strconv.Itoa(rm.retryAfter))
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		// wait or run
		span.AddEvent("queuing request")
		select {
		case rm.limiter <- struct{}{}:
			rm.runReqWithSync(c, span, reqLg)
		case <-c.Request.Context().Done():
			// reject with timeout
			rm.timedOut.Add(1)
			span.AddEvent("request's context expired before request was handled")
			reqLg.Error("request's context expired before request was handled")
			c.Header("Retry-After", strconv.Itoa(rm.retryAfter))
			c.AbortWithStatus(http.StatusTooManyRequests)
		}
	}
}

func (rm *RateLimiter) GetRunningRequests() int {
	return int(rm.running.Load())
}

func (rm *RateLimiter) GetTotalRequests() int {
	return int(rm.total.Load())
}

func (rm *RateLimiter) GetRejectedRequests() int {
	return int(rm.rejected.Load())
}

func (rm *RateLimiter) GetTimedOutRequests() int {
	return int(rm.timedOut.Load())
}

func (rm *RateLimiter) runReqWithSync(c *gin.Context, span trace.Span, reqLg *slog.Logger) {
	rm.running.Add(1)
	defer rm.running.Add(-1)
	defer func() { <-rm.limiter }()
	start := time.Now()
	span.AddEvent("request accepted")
	reqLg.Info("request accepted")
	c.Next()
	duration := time.Since(start)
	span.AddEvent(
		"request completed",
		trace.WithAttributes(
			attribute.Int("Status", c.Writer.Status()),
			attribute.String("Duration", duration.String()),
		),
	)
	reqLg.Info("request completed",
		"Status", c.Writer.Status(),
		"Duration", duration,
	)
}

func normalizeParams(maxRunning, maxWait, retryAfter int, lg *slog.Logger) (int, int, int) {
	if maxRunning < 1 {
		maxRunning = defaultMaxRunning
		lg.Warn("invalid max connection limit: limit was reset to default",
			"supplied limit", maxRunning,
			"default", defaultMaxRunning,
		)
	}

	if maxWait < 1 {
		maxWait = defaultMaxWait
		lg.Warn("invalid max total limit: limit was reset to default",
			"supplied limit", maxWait,
			"default", defaultMaxWait,
		)
	}

	if retryAfter < 1 {
		retryAfter = defaultRetryAfter
		lg.Warn("invalid max retry after: retry after was reset to default",
			"supplied limit", retryAfter,
			"default", defaultRetryAfter,
		)
	}
	return maxRunning, maxWait, retryAfter
}
