package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ec-recommend/internal/config"
	"ec-recommend/internal/handler"
	"ec-recommend/internal/repository"
	"ec-recommend/internal/router"
	"ec-recommend/internal/service"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize AWS configuration
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// Create database connection
	dbConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Set the global database for SQLBoiler
	boil.SetDB(db)

	// Create Bedrock runtime client
	bedrockClient := bedrockruntime.NewFromConfig(awsCfg)

	// Initialize services
	bedrockService := service.NewBedrockClient(bedrockClient, cfg.BedrockModelID)

	// Initialize repositories
	recommendationRepo := repository.NewRecommendationRepository(db)

	// Initialize recommendation service
	recommendationService := service.NewRecommendationService(recommendationRepo, bedrockService, cfg.BedrockModelID)

	// Initialize handlers
	aiHandler := handler.NewAIHandler(bedrockService)
	healthHandler := handler.NewHealthHandler()
	recommendationHandler := handler.NewRecommendationHandler(recommendationService)

	// Setup router
	routerEngine := router.SetupRouter(aiHandler, healthHandler, recommendationHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: routerEngine,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		log.Printf("Using AWS region: %s", cfg.AWSRegion)
		log.Printf("Using Bedrock model: %s", cfg.BedrockModelID)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
