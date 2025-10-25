package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"orderservice/internal/controller/http/handlers"
	"orderservice/internal/controller/http/middleware"
	"orderservice/internal/usecase"

	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger, u usecase.OrderUsecase) http.Handler {
	h := handlers.NewHandlers(logger, u)

	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(logger))

	r.Get("/", h.Order.Root)
	r.Get("/orders/{id}", h.Order.GetOrder)

	return r
}

type Server interface {
	Start()
	Shutdown(ctx context.Context) error
}

type serverImpl struct {
	srv    *http.Server
	logger *zap.Logger
}

func NewServer(logger *zap.Logger, handler http.Handler, addr string) Server {
	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return &serverImpl{srv: s, logger: logger}
}

func (s *serverImpl) Start() {
	go func() {
		s.logger.Info("HTTP server started", zap.String("address", s.srv.Addr))
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal("HTTP server error", zap.Error(err))
		}
	}()
}

func (s *serverImpl) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
