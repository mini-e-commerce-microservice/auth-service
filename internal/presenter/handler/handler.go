package handler

import (
	whttp "github.com/SyaibanAhmadRamadhan/http-wrapper"
	"github.com/go-chi/chi/v5"
	"github.com/mini-e-commerce-microservice/auth-service/internal/services/auth"
)

type handler struct {
	r           *chi.Mux
	authService auth.Service
	httpOtel    *whttp.Opentelemetry
}

func NewHandler(r *chi.Mux, authService auth.Service) {

	h := &handler{
		r:           r,
		authService: authService,
		httpOtel: whttp.NewOtel(
			whttp.WithRecoverMode(true),
			whttp.WithPropagator(),
			whttp.WithValidator(nil, nil),
		),
	}
	h.route()
}

func (h *handler) route() {
	h.r.Post("/v1/login", h.httpOtel.Trace(
		h.V1LoginPost, whttp.WithLogRequestBody(false), whttp.WithLogResponseBody(false),
	))

	h.r.Get("/v1/generate-access-token", h.httpOtel.Trace(
		h.V1GenerateAccessTokenPost, whttp.WithLogRequestBody(false),
	))

	h.r.Post("/v1/logout", h.httpOtel.Trace(
		h.V1LogoutPost, whttp.WithLogResponseBody(false), whttp.WithLogRequestBody(false),
	))
}
