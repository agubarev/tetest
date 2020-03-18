package server

import (
	"context"
	"net/http"

	"github.com/agubarev/tetest/internal/currency"
	"github.com/agubarev/tetest/internal/server/endpoints"
	"github.com/go-chi/chi"
)

type Server struct {
	manager *currency.Manager
}

func Run(ctx context.Context, m *currency.Manager, addr string) (err error) {
	r := chi.NewRouter()

	// route configuration
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/currency", func(r chi.Router) {
			r.Method("GET", "/", endpoints.NewEndpoint(m, endpoints.CurrencyGetLatest))
			r.Method("GET", "/{id}", endpoints.NewEndpoint(m, endpoints.CurrencyGetByID))
		})
	})

	return http.ListenAndServe(addr, r)
}
