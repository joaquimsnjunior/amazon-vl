package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"amazon-vl/internal/auth"

	goauth "github.com/abbot/go-http-auth"
)

// Config holds server configuration
type Config struct {
	Dir  string
	Port string
	Auth auth.Config
}

// Server represents the HTTP server
type Server struct {
	config     Config
	httpServer *http.Server
	auth       *goauth.BasicAuth
	fileServer *FileServer
}

// New creates a new Server with the given configuration
func New(cfg Config) *Server {
	return &Server{
		config:     cfg,
		auth:       auth.NewAuthenticator(cfg.Auth),
		fileServer: NewFileServer(cfg.Dir),
	}
}

// Run starts the server and handles graceful shutdown
func (s *Server) Run() error {
	// Setup HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.auth.Wrap(s.fileServer.Handle))
	mux.HandleFunc("/healthz", s.healthHandler)

	s.httpServer = &http.Server{
		Addr:         ":" + s.config.Port,
		Handler:      s.loggingMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for errors from server
	serverErrors := make(chan error, 1)

	// Start server
	go func() {
		log.Printf("INFO: Server starting on port %s", s.config.Port)
		log.Printf("INFO: Serving files from: %s", s.config.Dir)
		serverErrors <- s.httpServer.ListenAndServe()
	}()

	// Channel to listen for interrupt signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until error or shutdown signal
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Printf("INFO: Shutdown signal received: %v", sig)
		return s.Shutdown()
	}
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("INFO: Initiating graceful shutdown...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		// Force shutdown
		s.httpServer.Close()
		return fmt.Errorf("could not stop server gracefully: %w", err)
	}

	log.Println("INFO: Server stopped gracefully")
	return nil
}

// healthHandler returns server health status
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

// loggingMiddleware logs incoming requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response wrapper to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		log.Printf("INFO: %s %s %s %d %v",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			rw.statusCode,
			time.Since(start),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
