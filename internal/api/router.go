package api

import (
	"net/http"

	"truebit-api/internal/client"
	"truebit-api/internal/service"
)

func NewRouter(ethClient *client.EthereumClient) *http.ServeMux {
	truebitService := service.NewTruebitService(ethClient)
	handler := NewHandler(truebitService)

	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("GET /truebit", handler.GetTruebitInfo)

	return mux
}
