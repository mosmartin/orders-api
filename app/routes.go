package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mosmartin/orders-api/handler"
	"github.com/mosmartin/orders-api/repository/order"
)

func (a *App) loadRoutes() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/orders", a.loadOrderRoutes)

	a.router = r
}

func (a *App) loadOrderRoutes(r chi.Router) {
	orderHandler := &handler.Order{
		Repo: &order.RedisRepository{
			Client: a.rdb,
		},
	}

	r.Post("/", orderHandler.Create)
	r.Get("/", orderHandler.List)
	r.Get("/{id}", orderHandler.GetByID)
	r.Put("/{id}", orderHandler.UpdateByID)
	r.Delete("/{id}", orderHandler.DeleteByID)
}
