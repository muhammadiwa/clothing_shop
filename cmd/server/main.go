package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "clothing-shop-api/internal/api/routes"
    "clothing-shop-api/internal/config"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("could not load config: %v", err)
    }

    // Initialize router
    router := mux.NewRouter()

    // Set up routes
    routes.SetupRoutes(router)

    // Start the server
    log.Printf("Starting server on port %s...", cfg.ServerPort)
    if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
        log.Fatalf("could not start server: %v", err)
    }
}