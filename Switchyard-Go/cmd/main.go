package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.SetDefault("PORT", "8080")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// TODO: initialize DB connection pool (pgx)
	// TODO: run migrations (golang-migrate)
	// TODO: wire repositories → services → handlers

	port := viper.GetString("PORT")
	log.Printf("starting switchyard-go on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
