package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pecatu-be/internal/config"
	deliveryHttp "pecatu-be/internal/delivery/http"
	"pecatu-be/internal/repository/postgres"
	"pecatu-be/internal/usecase"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Setup Database
	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	// 3. Setup Repository
	donateRepo := postgres.NewDonateRepository(db)

	// 4. Setup Usecase
	donateUsecase := usecase.NewDonateUsecase(donateRepo)

	// 5. Setup Delivery (Handlers)
	healthHandler := deliveryHttp.NewHealthHandler()
	donateHandler := deliveryHttp.NewDonateHandler(donateUsecase)

	// 6. Setup Router
	router := deliveryHttp.NewRouter(healthHandler, donateHandler)

	// 7. Start Server with Graceful Shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("Server is starting on port %s...", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen and serve error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
