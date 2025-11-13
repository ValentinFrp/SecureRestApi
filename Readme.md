# Secure REST API - Clean Architecture

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://img.shields.io/badge/CI-Passing-success)](https://github.com/valentinfrappart/SecureRestApi/actions)
[![Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture-blueviolet)](ARCHITECTURE.md)

Une API REST sécurisée d'authentification en Go, implémentant la **Clean Architecture** (Ports & Adapters).

## 💡 Pourquoi ce projet ?

Ce projet démontre :
- **Clean Architecture** appliquée à Go (Ports & Adapters)
- **Sécurité** : JWT + bcrypt + bonnes pratiques
- **Testabilité** : Tests unitaires avec mocks, sans dépendances externes
- **SOLID Principles** : Dependency Inversion, Single Responsibility
- **Production-ready** : Docker, CI/CD, documentation complète

## 🛠️ Stack Technique

- **Language**: Go 1.21+
- **Architecture**: Clean Architecture (Ports & Adapters)
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **Database**: SQLite 3 (mattn/go-sqlite3)
- **Testing**: Go native testing + mocks
- **Containerization**: Docker (multi-stage builds)

## Architecture

Ce projet suit les principes de la **Clean Architecture** avec une séparation stricte des responsabilités :

```
SecureRestApi/
├── cmd/
│   └── api/
│       └── main.go                    # Point d'entrée de l'application
├── internal/
│   ├── domain/                        # Couche Domain (entités & règles métier)
│   │   ├── user.go                    # Entité User + interface Repository
│   │   └── errors.go                  # Erreurs métier
│   ├── usecase/                       # Couche Use Cases (logique applicative)
│   │   ├── auth_usecase.go            # Cas d'usage d'authentification
│   │   └── auth_usecase_test.go       # Tests unitaires
│   ├── infrastructure/                # Couche Infrastructure (implémentations)
│   │   ├── repository/
│   │   │   └── sqlite_user_repository.go  # Implémentation SQLite du UserRepository
│   │   ├── security/
│   │   │   ├── jwt.go                 # Service JWT
│   │   │   └── password.go            # Service de hashing bcrypt
│   │   └── database/
│   │       └── sqlite.go              # Connexion SQLite
│   └── delivery/                      # Couche Delivery (HTTP handlers)
│       └── http/
│           ├── handler.go             # Handlers des endpoints
│           ├── middleware.go          # Middlewares (auth, CORS, logs)
│           └── router.go              # Configuration des routes
├── go.mod
├── Dockerfile
└── README.md
```

### Avantages de la Clean Architecture

1. **Indépendance des frameworks** : Le domaine ne dépend pas des frameworks externes
2. **Testabilité** : Chaque couche peut être testée indépendamment
3. **Indépendance de la base de données** : Facile de changer SQLite pour PostgreSQL
4. **Maintenabilité** : Code organisé et prévisible
5. **Règle de dépendance** : Les dépendances pointent vers l'intérieur (domain ← usecase ← infrastructure/delivery)

## Fonctionnalités

- ✅ **Inscription** : Création de compte avec email/password
- ✅ **Connexion** : Authentification avec JWT
- ✅ **Route protégée** : Récupération du profil utilisateur authentifié
- ✅ **Sécurité** : Hashing bcrypt + JWT avec expiration
- ✅ **Base de données** : SQLite avec migrations automatiques
- ✅ **Tests unitaires** : Couverture des use cases

## Prérequis

- Go 1.21+
- Docker (optionnel)

## Installation & Exécution

### Option 1 : Avec Go

```bash
# Cloner le projet
git clone <repo>
cd SecureRestApi

# Installer les dépendances
go mod download

# Créer le dossier pour la base de données
mkdir -p data

# Lancer l'application
go run cmd/api/main.go
```

### Option 2 : Avec Docker

```bash
# Build l'image Docker
docker build -t secure-rest-api .

# Lancer le conteneur
docker run -p 8080:8080 -v $(pwd)/data:/root/data secure-rest-api
```

### Option 3 : Avec Docker Compose (recommandé)

```bash
# Lancer l'application
docker-compose up -d

# Voir les logs
docker-compose logs -f

# Arrêter l'application
docker-compose down
```

L'API sera disponible sur `http://localhost:8080`

## Tests

```bash
# Lancer tous les tests
go test ./...

# Tests avec couverture
go test -cover ./internal/usecase/

# Tests verbose
go test -v ./internal/usecase/
```

## Endpoints

### 1. Health Check (Public)
```bash
GET /health
```

**Réponse :**
```json
{
  "status": "healthy"
}
```

### 2. Inscription (Public)
```bash
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Réponse (201) :**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### 3. Connexion (Public)
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Réponse (200) :**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### 4. Profil Utilisateur (Protégé)
```bash
GET /api/auth/me
Authorization: Bearer <token>
```

**Réponse (200) :**
```json
{
  "id": 1,
  "email": "user@example.com",
  "created_at": "2024-01-15T10:30:00Z"
}
```

## Exemples Curl

```bash
# 1. Health check
curl http://localhost:8080/health

# 2. Inscription
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# 3. Connexion
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# 4. Profil (remplacer <TOKEN> par le token reçu)
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer <TOKEN>"
```

## Variables d'Environnement

| Variable | Description | Défaut |
|----------|-------------|---------|
| `PORT` | Port d'écoute du serveur | `8080` |
| `DB_PATH` | Chemin de la base SQLite | `./data/app.db` |
| `JWT_SECRET` | Clé secrète pour signer les JWT | `your-super-secret-key-change-this-in-production` |
| `JWT_ISSUER` | Émetteur du JWT | `secure-rest-api` |

**IMPORTANT** : En production, changez `JWT_SECRET` !

```bash
export JWT_SECRET="votre-cle-secrete-super-longue-et-aleatoire"
go run cmd/api/main.go
```

## Flux de données (Clean Architecture)

```
HTTP Request
    ↓
[Delivery Layer] handler.go → Reçoit la requête HTTP
    ↓
[Use Case Layer] auth_usecase.go → Exécute la logique métier
    ↓
[Domain Layer] user.go → Définit les règles métier
    ↓
[Infrastructure Layer] sqlite_user_repository.go → Persiste les données
    ↓
[Infrastructure Layer] security/jwt.go, password.go → Services techniques
```

## Concepts Implémentés

- **Ports & Adapters** : `UserRepository` est un port (interface), `SQLiteUserRepository` est un adaptateur
- **Dependency Injection** : Les dépendances sont injectées via les constructeurs
- **Separation of Concerns** : Chaque couche a une responsabilité unique
- **SOLID Principles** : Notamment le Dependency Inversion Principle
- **Test Doubles** : Mock repository pour tester les use cases en isolation

## Build de production

```bash
# Build binaire optimisé
CGO_ENABLED=1 go build -ldflags="-s -w" -o bin/api cmd/api/main.go

# Lancer le binaire
./bin/api
```

## Débogage

Les logs apparaissent dans stdout :
```
2024/01/15 10:30:00 Initializing database...
2024/01/15 10:30:00 Database initialized successfully
2024/01/15 10:30:00 🚀 Server starting on port 8080
```

## License

MIT
