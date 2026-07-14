package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crm-platform/backend/internal/auth"
	"github.com/crm-platform/backend/internal/crm"
	"github.com/crm-platform/backend/internal/pbx"
	"github.com/crm-platform/backend/internal/storage"
	"github.com/crm-platform/backend/internal/shared/middleware"
	"github.com/crm-platform/backend/pkg/broker"
	"github.com/crm-platform/backend/pkg/cache"
	"github.com/crm-platform/backend/pkg/database"
	"github.com/crm-platform/backend/pkg/logger"
	"github.com/crm-platform/backend/pkg/ws"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// ─── Initialize Logger ───────────────────────────────────
	log := logger.New(getEnv("APP_ENV", "development"))
	slog.SetDefault(log)

	slog.Info("🚀 Starting CRM & PBX Platform",
		"env", getEnv("APP_ENV", "development"),
		"port", getEnv("APP_PORT", "8080"),
	)

	// ─── Connect to PostgreSQL ───────────────────────────────
	db, err := database.NewPool(context.Background(), getEnv("DATABASE_URL",
		"postgres://crm_admin:change-me@localhost:5432/crm_platform?sslmode=disable"))
	if err != nil {
		slog.Error("Failed to connect to PostgreSQL", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("✅ Connected to PostgreSQL")

	// ─── Run Migrations ──────────────────────────────────────
	if err := database.RunMigrations(getEnv("DATABASE_URL",
		"postgres://crm_admin:change-me@localhost:5432/crm_platform?sslmode=disable"),
		"./migrations"); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("✅ Database migrations applied")

	// ─── Connect to Redis ────────────────────────────────────
	redisClient, err := cache.NewRedisClient(
		getEnv("REDIS_HOST", "localhost"),
		getEnv("REDIS_PORT", "6379"),
		getEnv("REDIS_PASSWORD", ""),
	)
	if err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()
	slog.Info("✅ Connected to Redis")

	// ─── Connect to RabbitMQ ─────────────────────────────────
	mq, err := broker.NewRabbitMQ(getEnv("RABBITMQ_URL",
		"amqp://crm_broker:change-me@localhost:5672/"))
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ", "error", err)
		os.Exit(1)
	}
	defer mq.Close()
	slog.Info("✅ Connected to RabbitMQ")

	// ─── Initialize WebSocket Hub ────────────────────────────
	wsHub := ws.NewHub(redisClient)
	go wsHub.Run()
	slog.Info("✅ WebSocket hub started")

	// ─── Initialize Modules ──────────────────────────────────
	// Auth module
	authService := auth.NewService(db, redisClient)
	authHandler := auth.NewHandler(authService)

	// CRM module
	crmRepo := crm.NewRepository(db)
	crmService := crm.NewService(crmRepo, mq, wsHub)
	crmHandler := crm.NewHandler(crmService)

	// Storage module
	storageService, err := storage.NewService(
		getEnv("MINIO_ENDPOINT", "localhost:9000"),
		getEnv("MINIO_ROOT_USER", "crm_storage_admin"),
		getEnv("MINIO_ROOT_PASSWORD", "change-me"),
		getEnv("MINIO_USE_SSL", "false") == "true",
	)
	if err != nil {
		slog.Error("Failed to connect to MinIO", "error", err)
		os.Exit(1)
	}
	storageHandler := storage.NewHandler(storageService)
	slog.Info("✅ Connected to MinIO")

	// PBX module
	pbxService := pbx.NewService(db, redisClient, mq, wsHub,
		getEnv("ASTERISK_ARI_URL", "http://localhost:8088/ari"),
		getEnv("ASTERISK_ARI_USER", "crm_ari"),
		getEnv("ASTERISK_ARI_PASSWORD", "change-me-ari-password"),
	)
	pbxHandler := pbx.NewHandler(pbxService)
	go pbxService.ConnectARI()            // Start ARI event loop
	go pbxService.StartRecordingConsumer() // Listen for recording upload events
	slog.Info("✅ PBX module initialized")

	// ─── Setup Router ────────────────────────────────────────
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.RequestID)
	r.Use(middleware.TraceID)
	r.Use(middleware.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(chimw.Compress(5))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{getEnv("FRONTEND_URL", "http://localhost:3000")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Trace-ID"},
		ExposedHeaders:   []string{"X-Trace-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"crm-platform"}`))
	})

	// Metrics endpoint (Prometheus)
	r.Get("/metrics", middleware.PrometheusHandler())

	// WebSocket endpoint
	r.Get("/ws", wsHub.HandleWebSocket)

	// ─── API Routes ──────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {
		// Auth (public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.RefreshToken)
			r.Post("/logout", authHandler.Logout)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(auth.JWTMiddleware(authService))
			r.Use(middleware.TenantContext)

			// CRM - Leads
			r.Route("/leads", func(r chi.Router) {
				r.Get("/", crmHandler.ListLeads)
				r.Post("/", crmHandler.CreateLead)
				r.Get("/{id}", crmHandler.GetLead)
				r.Put("/{id}", crmHandler.UpdateLead)
				r.Delete("/{id}", crmHandler.DeleteLead)
				r.Post("/{id}/convert", crmHandler.ConvertLead)
			})

			// CRM - Contacts
			r.Route("/contacts", func(r chi.Router) {
				r.Get("/", crmHandler.ListContacts)
				r.Post("/", crmHandler.CreateContact)
				r.Get("/{id}", crmHandler.GetContact)
				r.Put("/{id}", crmHandler.UpdateContact)
				r.Delete("/{id}", crmHandler.DeleteContact)
			})

			// CRM - Companies
			r.Route("/companies", func(r chi.Router) {
				r.Get("/", crmHandler.ListCompanies)
				r.Post("/", crmHandler.CreateCompany)
				r.Get("/{id}", crmHandler.GetCompany)
				r.Put("/{id}", crmHandler.UpdateCompany)
				r.Delete("/{id}", crmHandler.DeleteCompany)
			})

			// CRM - Deals
			r.Route("/deals", func(r chi.Router) {
				r.Get("/", crmHandler.ListDeals)
				r.Post("/", crmHandler.CreateDeal)
				r.Get("/{id}", crmHandler.GetDeal)
				r.Put("/{id}", crmHandler.UpdateDeal)
				r.Delete("/{id}", crmHandler.DeleteDeal)
			})

			// CRM - Pipelines
			r.Route("/pipelines", func(r chi.Router) {
				r.Get("/", crmHandler.ListPipelines)
				r.Post("/", crmHandler.CreatePipeline)
				r.Put("/{id}", crmHandler.UpdatePipeline)
			})

			// CRM - Activities
			r.Route("/activities", func(r chi.Router) {
				r.Get("/", crmHandler.ListActivities)
				r.Post("/", crmHandler.CreateActivity)
			})

			// PBX - Calls
			r.Route("/calls", func(r chi.Router) {
				r.Get("/active", pbxHandler.ListActiveCalls)
				r.Get("/history", pbxHandler.ListCallHistory)
				r.Post("/originate", pbxHandler.OriginateCall)
				r.Post("/{id}/transfer", pbxHandler.TransferCall)
				r.Post("/{id}/hold", pbxHandler.HoldCall)
				r.Post("/{id}/hangup", pbxHandler.HangupCall)
			})

			// PBX - Extensions
			r.Route("/pbx/extensions", func(r chi.Router) {
				r.Get("/", pbxHandler.ListExtensions)
				r.Post("/", pbxHandler.CreateExtension)
				r.Put("/{id}", pbxHandler.UpdateExtension)
				r.Delete("/{id}", pbxHandler.DeleteExtension)
			})

			// PBX - Routing
			r.Route("/pbx/routing", func(r chi.Router) {
				r.Get("/", pbxHandler.ListRoutingRules)
				r.Put("/", pbxHandler.UpdateRoutingRules)
			})

			// File storage
			r.Route("/files", func(r chi.Router) {
				r.Post("/upload", storageHandler.GetUploadURL)
				r.Get("/{id}", storageHandler.GetDownloadURL)
				r.Delete("/{id}", storageHandler.DeleteFile)
			})

			// Audit logs
			r.Get("/audit-logs", crmHandler.ListAuditLogs)
		})
	})

	// ─── Start Server ────────────────────────────────────────
	port := getEnv("APP_PORT", "8080")
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		slog.Info("🌐 Server listening", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("🛑 Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}
	slog.Info("👋 Server stopped")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
