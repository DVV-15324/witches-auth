// utils/blacklist_service.go
package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type BlacklistService struct {
	redis      *redis.Client
	revokedTTL int64
}

func NewBlacklistService(redis *redis.Client, revokedTTL int64) *BlacklistService {
	return &BlacklistService{redis: redis, revokedTTL: revokedTTL}
}

func (s *BlacklistService) cacheKeyBlacklist(accessToken string) string {
	return fmt.Sprintf("blacklist:%s", accessToken)
}

func (s *BlacklistService) BlacklistToken(ctx context.Context, accessToken string) error {
	if accessToken == "" {
		return fmt.Errorf("Error: accessToken cannot be empty")
	}

	key := s.cacheKeyBlacklist(accessToken)
	return s.redis.Set(ctx, key, "revoked", time.Duration(s.revokedTTL)*time.Second).Err()
}

// Việc Check nên để ở middleware
func (s *BlacklistService) IsTokenBlacklisted(ctx context.Context, token string) bool {
	if token == "" {
		return false
	}

	key := s.cacheKeyBlacklist(token)
	exists, err := s.redis.Exists(ctx, key).Result()
	return err == nil && exists == 1
}
