package http

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"orderservice/internal/controller/http/handlers"
	"orderservice/internal/usecase"
)

func NewRouter(u usecase.OrderUsecase) http.Handler {
	h := handlers.NewHandlers(u)

	r := chi.NewRouter()

	r.Get("/", h.Order.Root)
	r.Get("/orders/{id}", h.Order.GetOrder)

	return r
}

type Server interface {
	Start()
	Shutdown(ctx context.Context) error
}

type serverImpl struct {
	srv *http.Server
}

func NewServer(handler http.Handler, addr string) Server {
	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return &serverImpl{srv: s}
}

func (s *serverImpl) Start() {
	go func() {
		log.Println("HTTP server started at", s.srv.Addr)
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server error: %v", err)
		}
	}()
}

func (s *serverImpl) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
