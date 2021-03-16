package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/pechorka/kv/internal/store"
	"github.com/pechorka/kv/internal/web/mid"
)

func API(log *log.Logger, store store.Store) http.Handler {
	router := chi.NewRouter()

	router.Use(mid.Logger(log))
	router.Use(mid.Recoverer(log))
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
	})
	router.Use(corsMiddleware.Handler)

	kvGroup := KVGroup{Store: store}

	router.Route("/api/v1/", func(rv1 chi.Router) {
		rv1.Route("/kv", func(rKV chi.Router) {
			rKV.Get("/{key}", kvGroup.Get)
			rKV.Get("/", kvGroup.List)
			rKV.Post("/", kvGroup.Set)
			rKV.Delete("/{key}", kvGroup.Delete)
		})
	})

	return router
}
