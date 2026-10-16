package repository

import (
	"context"
	"encoding/json"
	"errors"
	modelRefresh "example/internal/refresh/model"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r *RefreshTokenRepository) cacheToken(ctx context.Context, token *modelRefresh.RefreshToken) error {
	if token == nil {
		return nil
	}
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	key := r.cacheKeyToken(token.Token)
	if err := r.core.Redis.GetClient().Set(ctx, key, data, 5*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RefreshTokenRepository) getCachedToken(ctx context.Context, token string) (*modelRefresh.RefreshToken, error) {
	key := r.cacheKeyToken(token)
	data, err := r.core.Redis.GetClient().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var rt modelRefresh.RefreshToken
	if err := json.Unmarshal([]byte(data), &rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *RefreshTokenRepository) deleteCacheToken(ctx context.Context, token string) error {
	key := r.cacheKeyToken(token)
	return r.core.Redis.GetClient().Del(ctx, key).Err()
}

func (r *RefreshTokenRepository) cacheJKT(ctx context.Context, token *modelRefresh.RefreshToken) error {
	if token == nil {
		return nil
	}
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	key := r.cacheKeyJKT(token.JKT)
	return r.core.Redis.GetClient().Set(ctx, key, data, 5*time.Minute).Err()
}

// getCachedJTK lấy refresh token theo JTK
func (r *RefreshTokenRepository) getCachedJKT(ctx context.Context, JKT string) (*modelRefresh.RefreshToken, error) {
	key := r.cacheKeyJKT(JKT)
	data, err := r.core.Redis.GetClient().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var rt modelRefresh.RefreshToken
	if err := json.Unmarshal([]byte(data), &rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

// deleteCacheJTK xóa cache theo JTK
func (r *RefreshTokenRepository) deleteCacheJTK(ctx context.Context, JKT string) error {
	key := r.cacheKeyJKT(JKT)
	return r.core.Redis.GetClient().Del(ctx, key).Err()
}
