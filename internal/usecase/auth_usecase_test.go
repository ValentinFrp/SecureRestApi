package usecase

import (
	"errors"
	"testing"

	"github.com/valentinfrappart/securerestapi/internal/domain"
	"github.com/valentinfrappart/securerestapi/internal/infrastructure/security"
)

type MockUserRepository struct {
	users         map[string]*domain.User
	nextID        int64
	createError   error
	findByIDError error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[string]*domain.User),
		nextID: 1,
	}
}

func (m *MockUserRepository) Create(email, passwordHash string) (*domain.User, error) {
	if m.createError != nil {
		return nil, m.createError
	}

	if _, exists := m.users[email]; exists {
		return nil, domain.ErrUserAlreadyExists
	}

	user := &domain.User{
		ID:           m.nextID,
		Email:        email,
		PasswordHash: passwordHash,
	}
	m.nextID++
	m.users[email] = user

	return user, nil
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) FindByID(id int64) (*domain.User, error) {
	if m.findByIDError != nil {
		return nil, m.findByIDError
	}

	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func TestAuthUseCase_Register_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	passwordService := security.NewPasswordService()
	jwtService := security.NewJWTService("test-secret", "test-issuer", 3600)

	useCase := NewAuthUseCase(mockRepo, passwordService, jwtService)

	req := RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := useCase.Register(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.Token == "" {
		t.Error("Expected token to be generated")
	}

	if resp.User.Email != req.Email {
		t.Errorf("Expected email %s, got %s", req.Email, resp.User.Email)
	}

	if resp.User.ID != 1 {
		t.Errorf("Expected ID 1, got %d", resp.User.ID)
	}
}

func TestAuthUseCase_Register_DuplicateEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	passwordService := security.NewPasswordService()
	jwtService := security.NewJWTService("test-secret", "test-issuer", 3600)

	useCase := NewAuthUseCase(mockRepo, passwordService, jwtService)

	req := RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	_, err := useCase.Register(req)
	if err != nil {
		t.Fatalf("Expected no error on first registration, got %v", err)
	}

	_, err = useCase.Register(req)
	if !errors.Is(err, domain.ErrUserAlreadyExists) {
		t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestAuthUseCase_Register_EmptyEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	passwordService := security.NewPasswordService()
	jwtService := security.NewJWTService("test-secret", "test-issuer", 3600)

	useCase := NewAuthUseCase(mockRepo, passwordService, jwtService)

	req := RegisterRequest{
		Email:    "",
		Password: "password123",
	}

	_, err := useCase.Register(req)
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthUseCase_Register_EmptyPassword(t *testing.T) {
	mockRepo := NewMockUserRepository()
	passwordService := security.NewPasswordService()
	jwtService := security.NewJWTService("test-secret", "test-issuer", 3600)

	useCase := NewAuthUseCase(mockRepo, passwordService, jwtService)

	req := RegisterRequest{
		Email:    "test@example.com",
		Password: "",
	}

	_, err := useCase.Register(req)
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthUseCase_Login_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	passwordService := security.NewPasswordService()
	jwtService := security.NewJWTService("test-secret", "test-issuer", 3600)

	useCase := NewAuthUseCase(mockRepo, passwordService, jwtService)

	registerReq := RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	_, err := useCase.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	loginReq := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := useCase.Login(loginReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.Token == "" {
		t.Error("Expected token to be generated")
	}

	if resp.User.Email != loginReq.Email {
		t.Errorf("Expected email %s, got %s", loginReq.Email, resp.User.Email)
	}
}

func TestAuthUseCase_Login_WrongPassword(t *testing.T) {
	mockRepo := NewMockUserRepository()
	passwordService := security.NewPasswordService()
	jwtService := security.NewJWTService("test-secret", "test-issuer", 3600)

	useCase := NewAuthUseCase(mockRepo, passwordService, jwtService)

	registerReq := RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	_, err := useCase.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	loginReq := LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, err = useCase.Login(loginReq)
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthUseCase_Login_UserNotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	passwordService := security.NewPasswordService()
	jwtService := security.NewJWTService("test-secret", "test-issuer", 3600)

	useCase := NewAuthUseCase(mockRepo, passwordService, jwtService)

	loginReq := LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	_, err := useCase.Login(loginReq)
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthUseCase_Login_EmptyEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	passwordService := security.NewPasswordService()
	jwtService := security.NewJWTService("test-secret", "test-issuer", 3600)

	useCase := NewAuthUseCase(mockRepo, passwordService, jwtService)

	loginReq := LoginRequest{
		Email:    "",
		Password: "password123",
	}

	_, err := useCase.Login(loginReq)
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthUseCase_GetUserByID_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	passwordService := security.NewPasswordService()
	jwtService := security.NewJWTService("test-secret", "test-issuer", 3600)

	useCase := NewAuthUseCase(mockRepo, passwordService, jwtService)

	registerReq := RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	resp, err := useCase.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	user, err := useCase.GetUserByID(resp.User.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Email != registerReq.Email {
		t.Errorf("Expected email %s, got %s", registerReq.Email, user.Email)
	}
}

func TestAuthUseCase_GetUserByID_NotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	passwordService := security.NewPasswordService()
	jwtService := security.NewJWTService("test-secret", "test-issuer", 3600)

	useCase := NewAuthUseCase(mockRepo, passwordService, jwtService)

	_, err := useCase.GetUserByID(999)
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}
