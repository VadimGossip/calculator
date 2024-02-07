package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type HttpServer struct {
	httpPort int
	*gin.Engine
	*http.Server
}

func NewHttpServer(httpPort int) *HttpServer {
	g := gin.New()
	g.RedirectTrailingSlash = false
	g.Use(gin.Recovery())
	s := &HttpServer{httpPort, g, nil}
	s.Server = &http.Server{
		Handler:        s,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s
}

func (s *HttpServer) Run() error {
	s.Addr = fmt.Sprintf(":%d", s.httpPort)
	logrus.Infof("Http server started at %d", s.httpPort)
	return s.Server.ListenAndServe()
}

func (s *HttpServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.Server.Shutdown(ctx)
}
