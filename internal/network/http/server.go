package http

import (
	"context"
	"github.com/KennyMacCormik/HerdMaster/internal/config"
	"github.com/KennyMacCormik/HerdMaster/internal/network"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	timeout time.Duration
	addr    string

	lg     *slog.Logger
	router *gin.Engine
	server *http.Server
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func New(conf config.Config, lg *slog.Logger) network.Endpoint {
	srv := Server{}

	srv.addr = strings.Join([]string{conf.Net.Host, strconv.Itoa(conf.Net.Port)}, ":")
	srv.timeout = conf.Net.Timeout
	srv.lg = lg

	srv.router = initGin(conf.Net.MaxConn)

	srv.server = &http.Server{
		Addr:         srv.addr,
		Handler:      srv.router,
		ReadTimeout:  srv.timeout,
		WriteTimeout: srv.timeout,
		IdleTimeout:  srv.timeout,
	}

	return &srv
}
