package handler

import (
	"net/http"
	"test_task/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/task/create", h.tCreate)
	mux.HandleFunc("/task/check-result", h.tStatus)

	return mux
}
