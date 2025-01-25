package middleware

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

func LogReq(c *gin.Context, uuid string, lg *slog.Logger, withReqDump bool) *slog.Logger {
	if withReqDump {
		lg.Debug("request received",
			"ClientIP", c.ClientIP(),
			"Proto", c.Request.Proto,
			"Header", c.Request.Header,
			"RemoteAddr", c.Request.RemoteAddr,
			"RequestURI", c.Request.RequestURI,
			"ContentLength", c.Request.ContentLength,
			"Method", c.Request.Method,
			"Host", c.Request.Host,
			"UrlPath", c.Request.URL.Path,
		)
	}
	return lg.With("UUID", uuid, "Method", c.Request.Method, "UrlPath", c.Request.URL.Path)
}
