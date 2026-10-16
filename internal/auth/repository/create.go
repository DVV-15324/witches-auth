package repository

import (
	"context"
	"example/internal/auth/mapping"
	domainAuth "example/internal/shared/domain"
)

func (r *AuthRepository) CreateAuth(ctx context.Context, domainAuth *domainAuth.Auth) error {
	auth := mapping.FromDomainToModelAuth(domainAuth)
	result := r.core.DB.WithContext(ctx).
		Select("email", "password", "user_id", "salt").
		Create(auth)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
