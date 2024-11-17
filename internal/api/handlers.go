package api

import (
	"encoding/json"
	"net/http"

	"truebit-api/internal/service"
)

type Handler struct {
	truebitService *service.TruebitService
}

func NewHandler(service *service.TruebitService) *Handler {
	return &Handler{
		truebitService: service,
	}
}

func (h *Handler) GetTruebitInfo(w http.ResponseWriter, r *http.Request) {
	info, err := h.truebitService.GetInfo(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info.Data)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(info.Display))
}
