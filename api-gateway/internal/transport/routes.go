// Package transport is a nice package
package transport

import "net/http"

func InitRoutes(h *Handler) *http.ServeMux {
	ordersMux := http.NewServeMux()
	
	ordersMux.HandleFunc("POST /orders", h.CreateOrder)
	ordersMux.HandleFunc("GET /orders/{id}", h.GetOrder) // GET /orders/{id}
	ordersMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
    h.log.Info().Msg("pong")
		w.Write([]byte("pong"))
	})

	return ordersMux
}
