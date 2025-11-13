package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	httpDelivery "github.com/valentinfrappart/securerestapi/internal/delivery/http"
	"github.com/valentinfrappart/securerestapi/internal/infrastructure/database"
	"github.com/valentinfrappart/securerestapi/internal/infrastructure/repository"
	"github.com/valentinfrappart/securerestapi/internal/infrastructure/security"
	"github.com/valentinfrappart/securerestapi/internal/usecase"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	}

	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "./data/app.db")
	jwtSecret := getEnv("JWT_SECRET", "your-super-secret-key-change-this-in-production")
	jwtIssuer := getEnv("JWT_ISSUER", "secure-rest-api")
	jwtDuration := 24 * time.Hour

	log.Println("Initializing database...")
	db, err := database.NewSQLiteDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	log.Println("Database initialized successfully")

	userRepo := repository.NewSQLiteUserRepository(db)
	passwordService := security.NewPasswordService()
	jwtService := security.NewJWTService(jwtSecret, jwtIssuer, jwtDuration)

	authUseCase := usecase.NewAuthUseCase(userRepo, passwordService, jwtService)

	handler := httpDelivery.NewHandler(authUseCase, jwtService)
	router := httpDelivery.NewRouter(handler, jwtService)

	mux := router.SetupRoutes()
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📚 API endpoints:")
	log.Printf("  - GET  /health              (public)")
	log.Printf("  - POST /api/auth/register  (public)")
	log.Printf("  - POST /api/auth/login     (public)")
	log.Printf("  - GET  /api/auth/me        (protected)")
	log.Println()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
