package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"cake3d_platform/internal/domain"
	"cake3d_platform/pkg/jwt"
	"cake3d_platform/pkg/password"
)

type userRepoStub struct {
	byEmail map[string]*domain.User
	byID    map[uint]*domain.User
	nextID  uint
}

func (r *userRepoStub) Create(_ context.Context, user *domain.User) error {
	r.nextID++
	user.ID = r.nextID
	r.byEmail[user.Email] = user
	r.byID[user.ID] = user
	return nil
}

func (r *userRepoStub) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	user, ok := r.byEmail[email]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *userRepoStub) GetByID(_ context.Context, id uint) (*domain.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

type tokenRepoStub struct {
	jti map[string]bool
}

func (r *tokenRepoStub) AddToBlacklist(_ context.Context, jti string, _ time.Duration) error {
	r.jti[jti] = true
	return nil
}

func (r *tokenRepoStub) IsBlacklisted(_ context.Context, jti string) (bool, error) {
	return r.jti[jti], nil
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := &userRepoStub{
		byEmail: map[string]*domain.User{
			"a@b.com": {ID: 1, Email: "a@b.com"},
		},
		byID:   map[uint]*domain.User{1: {ID: 1, Email: "a@b.com"}},
		nextID: 1,
	}
	tokenRepo := &tokenRepoStub{jti: map[string]bool{}}
	uc := NewAuthUseCase(repo, tokenRepo, jwt.NewManager("secret", "issuer", time.Hour))

	_, err := uc.Register(context.Background(), "test", "a@b.com", "123456")
	if !errors.Is(err, domain.ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists, got: %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	hash, _ := password.Hash("123456")
	user := &domain.User{ID: 1, Email: "a@b.com", PasswordHash: hash}
	repo := &userRepoStub{
		byEmail: map[string]*domain.User{"a@b.com": user},
		byID:    map[uint]*domain.User{1: user},
		nextID:  1,
	}
	tokenRepo := &tokenRepoStub{jti: map[string]bool{}}
	uc := NewAuthUseCase(repo, tokenRepo, jwt.NewManager("secret", "issuer", time.Hour))

	token, claims, err := uc.Login(context.Background(), "a@b.com", "123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("expected token")
	}
	if claims == nil || claims.UserID != 1 {
		t.Fatal("expected valid claims")
	}
}
