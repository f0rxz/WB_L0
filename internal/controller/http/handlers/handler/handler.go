package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"orderservice/internal/usecase"
)

type Handler struct {
	uc usecase.OrderUsecase
}

func NewHandler(u usecase.OrderUsecase) *Handler {
	return &Handler{uc: u}
}

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("./client/client.html")
	if err != nil {
		http.Error(w, "failed to read client.html", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write(data)
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")

	order, err := h.uc.GetOrder(r.Context(), orderID)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
