package redis

import (
	"context"
	"time"

	"cake3d_platform/internal/domain"
	goredis "github.com/go-redis/redis/v8"
)

const minBlacklistTTL = time.Second

type tokenRepository struct {
	client *goredis.Client
}

func NewTokenRepository(client *goredis.Client) domain.TokenRepository {
	return &tokenRepository{client: client}
}

func (r *tokenRepository) AddToBlacklist(ctx context.Context, jti string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = minBlacklistTTL
	}
	return r.client.Set(ctx, blacklistKey(jti), "1", ttl).Err()
}

func (r *tokenRepository) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	val, err := r.client.Exists(ctx, blacklistKey(jti)).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

func blacklistKey(jti string) string {
	return "auth:blacklist:" + jti
}
