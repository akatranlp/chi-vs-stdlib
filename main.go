package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

var useChi = false

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "chi" {
			useChi = true
		}
	}

	r := rootRoutes()
	http.ListenAndServe(":8080", LoggingMiddleware("root")(r))
}

func rootRoutes() http.Handler {
	var h http.Handler
	v1Mux := v1Routes()
	if useChi {
		log.Println("Using chi")
		r := chi.NewRouter()
		r.Mount("/api/v1", v1Mux)
		h = r
	} else {
		log.Println("Using stdlib")
		mux := http.NewServeMux()
		mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1Mux))
		// h = AppendSlashMiddleware(mux)
		h = RedirectSlashMiddleware(mux)
	}
	return h
}
func v1Routes() http.Handler {
	var h http.Handler
	todosMux := todoRoutes()
	if useChi {
		r := chi.NewRouter()
		r.Mount("/todos", todosMux)
		h = r
	} else {
		mux := http.NewServeMux()
		mux.Handle("/todos/", http.StripPrefix("/todos", todosMux))
		h = mux
	}
	return LoggingMiddleware("v1")(h)
}

func todoRoutes() http.Handler {
	var h http.Handler
	if useChi {
		r := chi.NewRouter()
		r.Get("/", HandleGetAllTodos)
		r.Post("/", HandleCreateTodo)
		r.Route("/{id}", func(r chi.Router) {
			r.Use(TodoMiddleware)
			r.Get("/", HandleGetTodo)
			r.Put("/", HandleUpdateTodo)
			r.Delete("/", HandleDeleteTodo)
		})
		h = r
	} else {
		mux := http.NewServeMux()

		mux.HandleFunc("GET /{$}", HandleGetAllTodos)
		mux.HandleFunc("POST /{$}", HandleCreateTodo)

		idMux := http.NewServeMux()
		idMux.HandleFunc("GET /{id}/{$}", HandleGetTodo)
		idMux.HandleFunc("PUT /{id}/{$}", HandleUpdateTodo)
		idMux.HandleFunc("DELETE /{id}/{$}", HandleDeleteTodo)

		mux.Handle("/{id}/{$}", TodoMiddleware(idMux))

		h = mux
	}
	return LoggingMiddleware("todos")(h)
}
