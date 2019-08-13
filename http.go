package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func newRouter() *chi.Mux {
	db := NewDB(os.Getenv("DNS_LEAK_REDIS_URI"))
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api", func(r chi.Router) {
		r.Get("/results/{id}", func(writer http.ResponseWriter, request *http.Request) {
			id := chi.URLParam(request, "id")

			results, err := db.GetResultsForId(id)
			if err != nil {
				r := ErrRender(err)
				_ = render.Render(writer, request, r)
				return
			}

			_ = render.Render(writer, request, results)
		})
	})

	r.Get("/*", func(writer http.ResponseWriter, request *http.Request) {
		// todo : load ui/build
	})

	return r
}
