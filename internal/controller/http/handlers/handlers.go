package handlers

import (
	"orderservice/internal/controller/http/handlers/handler"
	"orderservice/internal/usecase"

	"go.uber.org/zap"
)

type Handlers struct {
	Order  *handler.Handler
	logger *zap.Logger
}

func NewHandlers(logger *zap.Logger, u usecase.OrderUsecase) *Handlers {
	return &Handlers{
		Order:  handler.NewHandler(u, logger),
		logger: logger,
	}
}
