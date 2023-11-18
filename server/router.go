package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newRouter() *chi.Mux {
	r := chi.NewRouter()
	h := newHttpHandler()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome product api"))
	})

	r.Get("/products", h.getProducts)
	r.Get("/products/{id}", h.getProduct)
	r.Post("/products", h.postProduct)

	r.Post("/products/increase_inventory", h.increaseInventory)
	r.Post("/products/decrease_inventory", h.decreaseInventory)
	r.Post("/products/set_expired", h.setExpired)

	return r
}
