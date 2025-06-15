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
	bedrockRepository "ec-recommend/internal/repository/bedrock"
	dbRepository "ec-recommend/internal/repository/db"
	"ec-recommend/internal/router"
	"ec-recommend/internal/service"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println(".env file not found or failed to load, proceeding with system environment variables")
	}
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

	// Create AWS service clients
	bedrockClient := bedrockruntime.NewFromConfig(awsCfg)
	bedrockAgentClient := bedrockagentruntime.NewFromConfig(awsCfg)

	// Initialize V1 repositories
	bedrockRepo := bedrockRepository.NewBedrockClient(bedrockClient, cfg.BedrockModelID)
	recommendationRepo := dbRepository.NewRecommendationRepository(db)

	// Initialize V1 recommendation service
	recommendationService := service.NewRecommendationService(recommendationRepo, bedrockRepo, cfg.BedrockModelID)

	// Initialize V2 repositories
	recommendationRepoV2 := dbRepository.NewRecommendationRepositoryV2(db)
	bedrockRepoV2 := bedrockRepository.NewBedrockKnowledgeBaseService(bedrockAgentClient, bedrockClient, cfg.KnowledgeBaseID, cfg.BedrockModelID, cfg.EmbeddingModelID)

	// Initialize V2 services (Enhanced RAG-based)
	recommendationServiceV2 := service.NewRecommendationServiceV2(recommendationRepoV2, bedrockRepoV2, bedrockRepo, cfg.BedrockModelID, cfg.KnowledgeBaseID, cfg.EmbeddingModelID)

	// Initialize handlers
	chatHandler := handler.NewChatHandler(bedrockRepo)
	healthHandler := handler.NewHealthHandler()
	recommendationHandler := handler.NewRecommendationHandler(recommendationService)
	recommendationHandlerV2 := handler.NewRecommendationHandlerV2(recommendationServiceV2)

	// Setup router
	routerEngine := router.SetupRouter(chatHandler, healthHandler, recommendationHandler, recommendationHandlerV2)

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
		log.Printf("Using embedding model: %s", cfg.EmbeddingModelID)

		// Log V2 feature availability
		if cfg.KnowledgeBaseID != "" {
			log.Printf("V2 Knowledge Base features enabled with ID: %s", cfg.KnowledgeBaseID)
		}
		if cfg.OpenSearchEndpoint != "" {
			log.Printf("V2 Vector search features enabled with OpenSearch")
		}

		log.Println("Available endpoints:")
		log.Println("  V1 API: /api/v1/recommendations")
		log.Println("  V2 API: /api/v2/recommendations (Enhanced RAG-based)")
		log.Println("  Health: /health")

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
