package repository

import (
	"context"
	"errors"
	mapping "github.com/DVV-15324/witches-auth/internal/refresh/mapping"
	modelRefresh "github.com/DVV-15324/witches-auth/internal/refresh/model"
	domainRefresh "github.com/DVV-15324/witches-auth/internal/shared/domain"
	"gorm.io/gorm"
	"time"
)

func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*domainRefresh.RefreshToken, error) {
	if token == "" {
		return nil, nil
	}
	// 1. Check cache — lỗi cache không quan trọng, fallback DB
	if cached, err := r.getCachedToken(ctx, token); err == nil && cached != nil {
		return mapping.FromModelToDomainRefresh(cached), nil
	}
	// 2. Query DB
	var rt modelRefresh.RefreshToken
	err := r.core.DB.WithContext(ctx).
		Where("token = ?", token).
		First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	// 3. Cache nếu còn sống — best effort, bỏ qua lỗi
	if !rt.Revoked && time.Now().Before(rt.ExpiresAt) {
		_ = r.cacheToken(ctx, &rt)
	}
	return mapping.FromModelToDomainRefresh(&rt), nil
}

func (r *RefreshTokenRepository) GetByJKT(ctx context.Context, jkt string) (*domainRefresh.RefreshToken, error) {
	if jkt == "" {
		return nil, nil
	}
	// 1. Check cache
	if cached, err := r.getCachedJKT(ctx, jkt); err == nil && cached != nil {
		return mapping.FromModelToDomainRefresh(cached), nil
	}
	// 2. Query DB — chỉ lấy token còn sống
	var rt modelRefresh.RefreshToken
	err := r.core.DB.WithContext(ctx).
		Where("jkt = ? AND revoked = ? AND expires_at > ?", jkt, false, time.Now()).
		Order("created_at DESC").
		First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	// 3. Cache — best effort
	_ = r.cacheToken(ctx, &rt)
	return mapping.FromModelToDomainRefresh(&rt), nil
}
