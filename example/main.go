package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/induzo/crud/mock"
	"github.com/induzo/crud/rest"
)

func main() {
	r := chi.NewRouter()

	// Subrouters:
	m := mock.NewMgr()
	r.Route("/e", func(r chi.Router) {
		r.Get("/", rest.GETListHandler(m))
		r.Post("/", rest.POSTHandler(m))
		r.Get("/{ID}", rest.GETHandler(m))
		r.Patch("/{ID}", rest.PATCHHandler(m))
		r.Put("/{ID}", rest.PUTHandler(m))
		r.Delete("/{ID}", rest.DELETEHandler(m))
	})
	srv := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 7200 * time.Second,
		IdleTimeout:  10 * time.Second,
		Addr:         ":8080",
		Handler:      r,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
