package handlers

import (
	"encoding/json"
	"net/http"
	"rest-websockets/server"
)

type HomeResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func HomeHanlder(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(HomeResponse{
			Message: "Welcome to Mario's Page",
			Status:  true,
		})
	}
}
