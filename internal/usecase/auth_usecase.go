package usecase

import (
	"github.com/valentinfrappart/securerestapi/internal/domain"
	"github.com/valentinfrappart/securerestapi/internal/infrastructure/security"
)

type AuthUseCase struct {
	userRepo        domain.UserRepository
	passwordService *security.PasswordService
	jwtService      *security.JWTService
}

func NewAuthUseCase(
	userRepo domain.UserRepository,
	passwordService *security.PasswordService,
	jwtService *security.JWTService,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

func (uc *AuthUseCase) Register(req RegisterRequest) (*AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, domain.ErrInvalidCredentials
	}

	hashedPassword, err := uc.passwordService.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.Create(req.Email, hashedPassword)
	if err != nil {
		return nil, err
	}

	token, err := uc.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (uc *AuthUseCase) Login(req LoginRequest) (*AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, domain.ErrInvalidCredentials
	}

	user, err := uc.userRepo.FindByEmail(req.Email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := uc.passwordService.Verify(user.PasswordHash, req.Password); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := uc.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (uc *AuthUseCase) GetUserByID(id int64) (*domain.User, error) {
	return uc.userRepo.FindByID(id)
}
