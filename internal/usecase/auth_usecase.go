package usecase

import (
	"context"
	"errors"
	"time"

	"cake3d_platform/internal/domain"
	"cake3d_platform/pkg/jwt"
	"cake3d_platform/pkg/password"
)

type AuthUseCase interface {
	Register(ctx context.Context, name, email, rawPassword string) (*domain.User, error)
	Login(ctx context.Context, email, rawPassword string) (string, *jwt.Claims, error)
	Logout(ctx context.Context, claims *jwt.Claims) error
	GetProfile(ctx context.Context, userID uint) (*domain.User, error)
}

type authUseCase struct {
	userRepo  domain.UserRepository
	tokenRepo domain.TokenRepository
	jwt       *jwt.Manager
}

func NewAuthUseCase(userRepo domain.UserRepository, tokenRepo domain.TokenRepository, jwt *jwt.Manager) AuthUseCase {
	return &authUseCase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwt:       jwt,
	}
}

func (u *authUseCase) Register(ctx context.Context, name, email, rawPassword string) (*domain.User, error) {
	_, err := u.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, domain.ErrEmailAlreadyExists
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	hash, err := password.Hash(rawPassword)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
	}
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *authUseCase) Login(ctx context.Context, email, rawPassword string) (string, *jwt.Claims, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", nil, domain.ErrInvalidCredential
		}
		return "", nil, err
	}

	if err := password.Compare(user.PasswordHash, rawPassword); err != nil {
		return "", nil, domain.ErrInvalidCredential
	}

	token, claims, err := u.jwt.Generate(user.ID, user.Email)
	if err != nil {
		return "", nil, err
	}
	return token, claims, nil
}

func (u *authUseCase) Logout(ctx context.Context, claims *jwt.Claims) error {
	if claims == nil || claims.ExpiresAt == nil || claims.ID == "" {
		return domain.ErrUnauthorized
	}
	ttl := time.Until(claims.ExpiresAt.Time)
	return u.tokenRepo.AddToBlacklist(ctx, claims.ID, ttl)
}

func (u *authUseCase) GetProfile(ctx context.Context, userID uint) (*domain.User, error) {
	return u.userRepo.GetByID(ctx, userID)
}
