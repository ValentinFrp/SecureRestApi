.PHONY: help build run test clean docker-build docker-run install

APP_NAME=secure-rest-api
BINARY=main
DB_PATH=./data/app.db

help:
	@echo "Commandes disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install:
	@echo "📦 Installation des dépendances..."
	go mod download
	go mod tidy
	@echo "✅ Dépendances installées"

build:
	@echo "🔨 Compilation..."
	@mkdir -p data
	CGO_ENABLED=1 go build -o $(BINARY) cmd/api/main.go
	@echo "✅ Compilé: ./$(BINARY)"

run:
	@echo "🚀 Démarrage de l'API..."
	@mkdir -p data
	go run cmd/api/main.go

test:
	@echo "🧪 Lancement des tests..."
	go test -v ./internal/usecase/

test-coverage:
	@echo "🧪 Lancement des tests avec couverture..."
	go test -cover ./internal/usecase/
	go test -coverprofile=coverage.out ./internal/usecase/
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Rapport de couverture généré: coverage.html"

test-api: build
	@echo "🧪 Tests d'intégration de l'API..."
	@echo "⚠️  Assurez-vous que l'API tourne sur le port 8080"
	@./test_api.sh

lint:
	@echo "🔍 Vérification du code..."
	go vet ./...
	go fmt ./...
	@echo "✅ Code vérifié"

clean:
	@echo "🧹 Nettoyage..."
	rm -f $(BINARY)
	rm -f coverage.out coverage.html
	rm -rf data/*.db
	@echo "✅ Nettoyage terminé"

docker-build:
	@echo "🐳 Build de l'image Docker..."
	docker build -t $(APP_NAME) .
	@echo "✅ Image Docker créée: $(APP_NAME)"

docker-run:
	@echo "🐳 Lancement du conteneur Docker..."
	docker run -p 8080:8080 -v $(PWD)/data:/root/data $(APP_NAME)

docker-stop:
	@echo "🛑 Arrêt des conteneurs..."
	docker stop $$(docker ps -q --filter ancestor=$(APP_NAME)) 2>/dev/null || true

dev:
	@which air > /dev/null || (echo "❌ 'air' n'est pas installé. Installez-le avec: go install github.com/cosmtrek/air@latest" && exit 1)
	air

db-reset:
	@echo "🗑️  Suppression de la base de données..."
	rm -f $(DB_PATH)
	@echo "✅ Base de données réinitialisée"

all: clean install build test
	@echo "✅ Tout est prêt!"
