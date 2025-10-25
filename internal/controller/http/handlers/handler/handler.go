package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"orderservice/internal/controller/http/middleware"
	"orderservice/internal/usecase"
)

type Handler struct {
	uc     usecase.OrderUsecase
	logger *zap.Logger
}

func NewHandler(u usecase.OrderUsecase, logger *zap.Logger) *Handler {
	return &Handler{uc: u, logger: logger}
}

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestID(r.Context())

	data, err := os.ReadFile("./client/client.html")
	if err != nil {
		h.logger.Error("failed to read client.html", zap.String("request_id", reqID), zap.Error(err))
		http.Error(w, "failed to read client.html", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Request-ID", reqID)
	w.WriteHeader(http.StatusOK)
	w.Write(data)

	h.logger.Info("served client.html", zap.String("request_id", reqID))
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestID(r.Context())
	orderID := chi.URLParam(r, "id")

	order, err := h.uc.GetOrder(r.Context(), orderID)
	if err != nil {
		h.logger.Warn("order not found",
			zap.String("request_id", reqID),
			zap.String("order_id", orderID),
			zap.Error(err),
		)
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", reqID)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		h.logger.Error("failed to encode response",
			zap.String("request_id", reqID),
			zap.Error(err),
		)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	h.logger.Info("order retrieved",
		zap.String("request_id", reqID),
		zap.String("order_id", orderID),
	)
}
