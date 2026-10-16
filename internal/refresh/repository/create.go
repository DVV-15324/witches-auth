package repository

import (
	"context"
	mapping "example/internal/refresh/mapping"
	domainRefresh "example/internal/shared/domain"
)

func (r *RefreshTokenRepository) Create(ctx context.Context, token *domainRefresh.RefreshToken) error {
	modelToken := mapping.FromDomainToModelRefresh(token)
	if err := r.core.DB.WithContext(ctx).Create(modelToken).Error; err != nil {
		return err
	}
	return r.cacheToken(ctx, modelToken)
}
