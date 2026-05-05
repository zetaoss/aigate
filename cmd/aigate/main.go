package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"aigate/internal/config"
	"aigate/internal/gateway"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to YAML gateway config")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	gw, err := gateway.NewGateway(cfg)
	if err != nil {
		log.Fatalf("failed to initialize gateway: %v", err)
	}

	gw.LoadAvailableModels(context.Background())
	if cfg.Server.ValidateModelsOnStartup != nil && *cfg.Server.ValidateModelsOnStartup {
		if err := gw.ValidateStartupModels(); err != nil {
			log.Fatalf("failed startup model validation: %v", err)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", gw.HealthHandler)
	mux.HandleFunc("/v1/models", gw.ModelsHandler)
	mux.HandleFunc("/v1/models/", gw.ModelsByProviderHandler)
	mux.HandleFunc("/v1/chat/completions", gw.ChatCompletionsHandler)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:           gateway.LoggingMiddleware(gw, mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("starting AI gateway on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
