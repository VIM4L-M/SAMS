package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/vimal/sams/config"
	"github.com/vimal/sams/internal/engine"
	"github.com/vimal/sams/internal/handlers"
	"github.com/vimal/sams/internal/middleware"
)

func main() {
	cfg := config.Load()

	rules, err := engine.LoadRules(cfg.RulesPath)
	if err != nil {
		log.Fatalf("failed to load rules from %s: %v", cfg.RulesPath, err)
	}
	log.Printf("loaded %d rules from %s", len(rules), cfg.RulesPath)

	validator := engine.NewValidator(rules)

	rl := middleware.NewRateLimiter(10, 30) // 10 req/s, burst of 30

	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.CORS(cfg.AllowedOrigins))
	r.Use(rl.Middleware)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/validate", handlers.Validate(validator))
		r.Get("/components", handlers.Components)
		r.Get("/rules", handlers.Rules(rules))
		r.Get("/health", handlers.Health(len(rules)))
	})

	addr := ":" + cfg.Port
	log.Printf("SAMS backend listening on %s (env=%s)", addr, cfg.Env)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
