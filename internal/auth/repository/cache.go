package repository

import (
	"context"
	"encoding/json"
	"errors"

	modelAuth "example/internal/auth/model"

	"time"

	"github.com/redis/go-redis/v9"
)

func (r *AuthRepository) cacheAuth(ctx context.Context, auth *modelAuth.Auth) error {
	if auth == nil {
		return nil
	}

	pipe := r.core.Redis.GetClient().Pipeline()

	// Cache by email
	emailKey := r.cacheKeyAuthByEmail(auth.Email)
	emailData, _ := json.Marshal(auth)
	pipe.Set(ctx, emailKey, emailData, 5*time.Minute)

	// Cache by user_id
	userKey := r.cacheKeyAuthByUserID(auth.UserId)
	userData, _ := json.Marshal(auth)
	pipe.Set(ctx, userKey, userData, 5*time.Minute)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *AuthRepository) getAuthFromCacheByEmail(ctx context.Context, email string) (*modelAuth.Auth, error) {
	key := r.cacheKeyAuthByEmail(email)
	data, err := r.core.Redis.GetClient().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var auth modelAuth.Auth
	if err := json.Unmarshal([]byte(data), &auth); err != nil {
		return nil, err
	}
	return &auth, nil
}

func (r *AuthRepository) getAuthFromCacheByUserID(ctx context.Context, userID int64) (*modelAuth.Auth, error) {
	key := r.cacheKeyAuthByUserID(userID)
	data, err := r.core.Redis.GetClient().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var auth modelAuth.Auth
	if err := json.Unmarshal([]byte(data), &auth); err != nil {
		return nil, err
	}
	return &auth, nil
}

func (r *AuthRepository) invalidateAuthCache(ctx context.Context, email string, userID int64) {
	r.core.Redis.GetClient().Del(ctx, r.cacheKeyAuthByEmail(email))
	r.core.Redis.GetClient().Del(ctx, r.cacheKeyAuthByUserID(userID))
}

func (r *AuthRepository) cacheNotFound(ctx context.Context, key string) {
	r.core.Redis.GetClient().Set(ctx, key, "1", 1*time.Minute)
}
