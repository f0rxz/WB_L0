package handlers

import (
	"orderservice/internal/controller/http/handlers/handler"
	"orderservice/internal/usecase"
)

type Handlers struct {
	Order *handler.Handler
}

func NewHandlers(u usecase.OrderUsecase) *Handlers {
	return &Handlers{
		Order: handler.NewHandler(u),
	}
}
