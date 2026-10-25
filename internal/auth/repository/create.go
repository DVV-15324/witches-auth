package repository

import (
	"context"
	"github.com/DVV-15324/witches-auth/internal/auth/mapping"
	domainAuth "github.com/DVV-15324/witches-auth/internal/shared/domain"
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
