package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors" // Add this import
	"github.com/shanurrahman/orchestrator/config"
	"github.com/shanurrahman/orchestrator/docker"
	"github.com/shanurrahman/orchestrator/docs"
	"github.com/shanurrahman/orchestrator/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Orchestrator API
// @version         1.0
// @description     A container orchestration service API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8090
// @BasePath  /
// @schemes   http https
func main() {
	log.Println("Starting the orchestrator service...")

	cfg := config.Load()
	log.Println("Configuration loaded successfully")

	// If running behind a proxy (e.g. for HTTPS termination), update Swagger settings
	if cfg.BehindProxy {
		// Set host to empty so Swagger uses the request's host
		docs.SwaggerInfo.Host = ""
		// If you’re using a base path behind the proxy, update it accordingly
		docs.SwaggerInfo.BasePath = "/orchestrator"
		// Force Swagger to use https
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		docs.SwaggerInfo.Host = "localhost:8090"
		docs.SwaggerInfo.BasePath = "/"
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	dockerClient := docker.NewDockerManager(cfg)
	log.Println("Docker manager initialized")

	// Create a new chi router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	
	// Add CORS middleware
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"*"}, // Allow all origins
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Health check endpoint
	healthHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}
	r.Get("/health", healthHandler)
	r.Head("/health", healthHandler)

	// Serve Swagger documentation
	swaggerURL := "/swagger/doc.json"
	// If behind a proxy, the doc.json might be served from a different path
	if cfg.BehindProxy {
		swaggerURL = "/orchestrator/swagger/doc.json"
	}
	r.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(swaggerURL),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.PersistAuthorization(true),
	))

	// Register API routes
	r.Route("/containers", func(r chi.Router) {
		r.Get("/images", handlers.ListImagesHandler(dockerClient))
		r.Post("/", handlers.CreateContainerHandler(dockerClient))
		r.Get("/{id}/status", handlers.GetContainerStatusHandler(dockerClient))
		r.Delete("/{id}/kill", handlers.KillContainerHandler(dockerClient))
	})

	// Configure server with timeouts
	server := &http.Server{
		Addr:         "0.0.0.0:8090",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Printf("Server starting on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}