# Architecture Documentation

## Clean Architecture - Vue d'ensemble

Ce projet implémente la **Clean Architecture** (aussi appelée **Hexagonal Architecture** ou **Ports & Adapters**), un pattern architectural qui favorise la séparation des préoccupations et l'indépendance vis-à-vis des frameworks et technologies externes.

## Principes Fondamentaux

### 1. Règle de Dépendance
Les dépendances pointent **toujours vers l'intérieur** :
```
Infrastructure/Delivery → Use Cases → Domain
```

Le **Domain** ne dépend de rien. Les **Use Cases** ne dépendent que du Domain. L'**Infrastructure** et la **Delivery** dépendent des Use Cases et du Domain.

### 2. Séparation des Couches

```
┌─────────────────────────────────────────────────────────────┐
│                      DELIVERY LAYER                          │
│  (HTTP Handlers, Middleware, Router)                        │
│  → Reçoit les requêtes, sérialise/désérialise               │
└──────────────────────┬──────────────────────────────────────┘
                       │ Appelle
                       ↓
┌─────────────────────────────────────────────────────────────┐
│                     USE CASE LAYER                           │
│  (Business Logic Application)                                │
│  → Orchestre les entités et services du domaine             │
└──────────────────────┬──────────────────────────────────────┘
                       │ Utilise
                       ↓
┌─────────────────────────────────────────────────────────────┐
│                      DOMAIN LAYER                            │
│  (Entités, Interfaces, Règles Métier)                       │
│  → Définit les règles métier pures                          │
└──────────────────────┬──────────────────────────────────────┘
                       ↑ Implémente
                       │
┌─────────────────────────────────────────────────────────────┐
│                  INFRASTRUCTURE LAYER                        │
│  (Database, Security, External Services)                     │
│  → Détails techniques d'implémentation                      │
└─────────────────────────────────────────────────────────────┘
```

## Structure des Couches

### Domain Layer (`internal/domain/`)

**Responsabilité** : Contient les entités et les règles métier pures.

**Fichiers** :
- `user.go` : Entité User + interface UserRepository (Port)
- `errors.go` : Erreurs métier

**Caractéristiques** :
- ❌ **Aucune dépendance externe** (pas d'import de packages tiers)
- ✅ Définit les **interfaces (ports)** que les autres couches doivent implémenter
- ✅ Contient la **logique métier pure**
- ✅ Indépendant de la base de données, du framework web, etc.

**Exemple - Le Port UserRepository** :
```go
type UserRepository interface {
    Create(email, passwordHash string) (*User, error)
    FindByEmail(email string) (*User, error)
    FindByID(id int64) (*User, error)
}
```
C'est une **interface** définie dans le domain. L'implémentation concrète est dans l'infrastructure.

---

### Use Case Layer (`internal/usecase/`)

**Responsabilité** : Contient la logique applicative (orchestration).

**Fichiers** :
- `auth_usecase.go` : Cas d'usage d'authentification (Register, Login, GetUserByID)
- `auth_usecase_test.go` : Tests unitaires avec mocks

**Caractéristiques** :
- ✅ Dépend **uniquement** du Domain
- ✅ Orchestre les entités et les services
- ✅ Appelle les repositories via les interfaces (Dependency Inversion)
- ✅ Facilement testable avec des mocks

**Exemple - Register Use Case** :
```go
func (uc *AuthUseCase) Register(req RegisterRequest) (*AuthResponse, error) {
    // 1. Validation
    if req.Email == "" || req.Password == "" {
        return nil, domain.ErrInvalidCredentials
    }
    
    // 2. Hash du mot de passe (via service d'infrastructure)
    hashedPassword, err := uc.passwordService.Hash(req.Password)
    
    // 3. Créer l'utilisateur (via le port UserRepository)
    user, err := uc.userRepo.Create(req.Email, hashedPassword)
    
    // 4. Générer le token JWT
    token, err := uc.jwtService.GenerateToken(user.ID, user.Email)
    
    return &AuthResponse{Token: token, User: user}, nil
}
```

---

### Infrastructure Layer (`internal/infrastructure/`)

**Responsabilité** : Implémente les détails techniques (adapters).

**Structure** :
```
infrastructure/
├── repository/
│   └── sqlite_user_repository.go  # Implémentation SQLite du UserRepository
├── security/
│   ├── jwt.go                     # Service JWT
│   └── password.go                # Service bcrypt
└── database/
    └── sqlite.go                  # Connexion SQLite
```

**Caractéristiques** :
- ✅ Implémente les **interfaces (ports)** définies dans le Domain
- ✅ Contient les détails techniques (SQL, crypto, etc.)
- ✅ **Remplaçable facilement** (ex: SQLite → PostgreSQL)

**Exemple - Adapter SQLiteUserRepository** :
```go
type SQLiteUserRepository struct {
    db *sql.DB
}

// Implémente l'interface domain.UserRepository
func (r *SQLiteUserRepository) Create(email, passwordHash string) (*domain.User, error) {
    // Détails SQL spécifiques à SQLite
    query := `INSERT INTO users (email, password_hash, ...) VALUES (?, ?, ...)`
    result, err := r.db.Exec(query, email, passwordHash, ...)
    // ...
}
```

---

### Delivery Layer (`internal/delivery/http/`)

**Responsabilité** : Gère la communication HTTP avec le monde extérieur.

**Fichiers** :
- `handler.go` : Handlers des endpoints (Register, Login, Me, Health)
- `middleware.go` : Middlewares (AuthMiddleware, CORS, Logging)
- `router.go` : Configuration des routes

**Caractéristiques** :
- ✅ Reçoit les requêtes HTTP
- ✅ Valide et désérialise les données
- ✅ Appelle les **Use Cases**
- ✅ Sérialise les réponses en JSON

**Exemple - Handler Register** :
```go
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    // 1. Désérialiser la requête HTTP
    var req usecase.RegisterRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // 2. Appeler le use case
    resp, err := h.authUseCase.Register(req)
    
    // 3. Sérialiser et renvoyer la réponse
    respondWithJSON(w, http.StatusCreated, resp)
}
```

---

## Flux de Données Complet

### Exemple : Inscription d'un utilisateur

```
1. Client HTTP
   POST /api/auth/register
   { "email": "user@example.com", "password": "pass123" }
         │
         ↓
2. [DELIVERY] handler.Register()
   - Désérialise le JSON
   - Crée RegisterRequest
         │
         ↓
3. [USE CASE] authUseCase.Register()
   - Valide les données
   - Hash le password via PasswordService
   - Appelle userRepo.Create()
         │
         ↓
4. [INFRASTRUCTURE] sqliteUserRepository.Create()
   - Exécute la requête SQL INSERT
   - Retourne *domain.User
         │
         ↓
5. [USE CASE] authUseCase.Register()
   - Génère le JWT via JWTService
   - Retourne AuthResponse
         │
         ↓
6. [DELIVERY] handler.Register()
   - Sérialise en JSON
   - Envoie la réponse HTTP 201
         │
         ↓
7. Client HTTP
   { "token": "eyJ...", "user": {...} }
```

---

## Ports & Adapters

### Qu'est-ce qu'un Port ?
Un **port** est une **interface** définie dans le Domain qui décrit un contrat.

**Exemple** :
```go
// Port défini dans domain/user.go
type UserRepository interface {
    Create(email, passwordHash string) (*User, error)
    FindByEmail(email string) (*User, error)
    FindByID(id int64) (*User, error)
}
```

### Qu'est-ce qu'un Adapter ?
Un **adapter** est une **implémentation concrète** d'un port.

**Exemple** :
```go
// Adapter SQLite dans infrastructure/repository/sqlite_user_repository.go
type SQLiteUserRepository struct {
    db *sql.DB
}

func (r *SQLiteUserRepository) Create(...) (*domain.User, error) {
    // Implémentation SQLite
}
```

### Pourquoi c'est puissant ?

Si demain vous voulez passer de SQLite à PostgreSQL :
1. Créez `PostgresUserRepository` qui implémente `domain.UserRepository`
2. Changez l'injection de dépendance dans `main.go`
3. **Aucun changement dans le Domain ou les Use Cases** ! 🎉

---

## Testabilité

La Clean Architecture rend les tests **extrêmement simples**.

### Tests Unitaires des Use Cases

Vous pouvez tester les Use Cases **sans base de données réelle** :

```go
// Mock du repository (dans auth_usecase_test.go)
type MockUserRepository struct {
    users map[string]*domain.User
}

func (m *MockUserRepository) Create(email, passwordHash string) (*domain.User, error) {
    if _, exists := m.users[email]; exists {
        return nil, domain.ErrUserAlreadyExists
    }
    user := &domain.User{ID: 1, Email: email, PasswordHash: passwordHash}
    m.users[email] = user
    return user, nil
}

// Test
func TestAuthUseCase_Register_Success(t *testing.T) {
    mockRepo := NewMockUserRepository()
    useCase := NewAuthUseCase(mockRepo, passwordService, jwtService)
    
    resp, err := useCase.Register(RegisterRequest{
        Email: "test@example.com",
        Password: "password123",
    })
    
    assert.NoError(t, err)
    assert.NotNil(t, resp.Token)
}
```

✅ **Pas de base de données**
✅ **Pas de serveur HTTP**
✅ **Tests ultra-rapides**

---

## Dependency Inversion Principle (SOLID)

La Clean Architecture respecte le **D** de SOLID.

### ❌ Sans Dependency Inversion (mauvais)
```
[Use Case] → dépend de → [SQLiteRepository]
```
Si vous changez de base de données, vous devez modifier le Use Case !

### Avec Dependency Inversion (bon)
```
[Use Case] → dépend de → [UserRepository Interface]
                              ↑
                              │ implémente
                              │
                      [SQLiteRepository]
```
Le Use Case dépend de l'**abstraction** (interface), pas de l'implémentation concrète.

---

## Avantages de cette Architecture

| Avantage | Description |
|----------|-------------|
| **Indépendance des frameworks** | Le domaine ne dépend pas de Gin, Echo, etc. |
| **Testabilité** | Chaque couche peut être testée isolément |
| **Maintenabilité** | Code organisé et prévisible |
| **Flexibilité** | Facile de changer de BDD, de framework, etc. |
| **Scalabilité** | Structure claire pour les grandes applications |
| **Onboarding** | Les nouveaux dev comprennent vite la structure |

---

## Concepts Clés Implémentés

### 1. Separation of Concerns
Chaque couche a **une seule responsabilité** :
- Domain : règles métier
- Use Case : orchestration
- Infrastructure : détails techniques
- Delivery : communication HTTP

### 2. Dependency Injection
Les dépendances sont injectées via les constructeurs :
```go
func NewAuthUseCase(
    userRepo domain.UserRepository,      // Interface, pas implémentation
    passwordService *security.PasswordService,
    jwtService *security.JWTService,
) *AuthUseCase {
    return &AuthUseCase{
        userRepo: userRepo,
        passwordService: passwordService,
        jwtService: jwtService,
    }
}
```

### 3. Interface Segregation
Les interfaces sont **petites et spécifiques** :
```go
type UserRepository interface {
    Create(email, passwordHash string) (*User, error)
    FindByEmail(email string) (*User, error)
    FindByID(id int64) (*User, error)
}
```

---

## Extension du Projet

### Ajouter une nouvelle fonctionnalité "Reset Password"

#### 1. Domain Layer
```go
// domain/user.go
type UserRepository interface {
    // ... méthodes existantes
    UpdatePassword(userID int64, newPasswordHash string) error
}
```

#### 2. Use Case Layer
```go
// usecase/auth_usecase.go
func (uc *AuthUseCase) ResetPassword(userID int64, newPassword string) error {
    hashedPassword, err := uc.passwordService.Hash(newPassword)
    if err != nil {
        return err
    }
    return uc.userRepo.UpdatePassword(userID, hashedPassword)
}
```

#### 3. Infrastructure Layer
```go
// infrastructure/repository/sqlite_user_repository.go
func (r *SQLiteUserRepository) UpdatePassword(userID int64, newPasswordHash string) error {
    query := `UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`
    _, err := r.db.Exec(query, newPasswordHash, time.Now(), userID)
    return err
}
```

#### 4. Delivery Layer
```go
// delivery/http/handler.go
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
    var req ResetPasswordRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    userID := r.Context().Value(contextKeyUserID).(int64)
    err := h.authUseCase.ResetPassword(userID, req.NewPassword)
    
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJSON(w, http.StatusOK, map[string]string{"message": "Password updated"})
}
```

---

## Ressources

- [The Clean Architecture (Robert C. Martin)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture (Alistair Cockburn)](https://alistair.cockburn.us/hexagonal-architecture/)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)

---

## Conclusion

La Clean Architecture peut sembler **over-engineered** pour un petit projet, mais elle brille dans les applications moyennes à grandes. Elle garantit :

✅ Un code **maintenable**
✅ Des tests **faciles et rapides**
✅ Une **flexibilité** pour changer de technologies
✅ Une **onboarding** simplifié pour les nouveaux développeurs

C'est un investissement initial qui paie sur le long terme ! 🚀
