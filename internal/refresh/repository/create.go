package repository

import (
	"context"
	mapping "github.com/DVV-15324/witches-auth/internal/refresh/mapping"
	domainRefresh "github.com/DVV-15324/witches-auth/internal/shared/domain"
)

func (r *RefreshTokenRepository) Create(ctx context.Context, token *domainRefresh.RefreshToken) error {
	modelToken := mapping.FromDomainToModelRefresh(token)
	if err := r.core.DB.WithContext(ctx).Create(modelToken).Error; err != nil {
		return err
	}
	return r.cacheToken(ctx, modelToken)
}
