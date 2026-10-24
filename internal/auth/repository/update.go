package repository

import (
	"context"

	"github.com/DVV-15324/witches-auth/internal/auth/mapping"
	modelAuth "github.com/DVV-15324/witches-auth/internal/auth/model"
	domainAuth "github.com/DVV-15324/witches-auth/pkg/shared/domain"
)

func (r *AuthRepository) UpdateAuth(ctx context.Context, modleAuth *domainAuth.Auth) error {
	auth := mapping.FromDomainToModelAuth(modleAuth)
	// Lấy old auth để invalidate cache
	var oldAuth modelAuth.Auth
	r.core.DB.WithContext(ctx).Where("email = ?", auth.Email).First(&oldAuth)

	result := r.core.DB.WithContext(ctx).Model(&modelAuth.Auth{}).
		Where("email = ?", auth.UserId).
		Updates(map[string]interface{}{
			"password":  auth.Password,
			"salt":      auth.Salt,
			"banned":    auth.Banned,
			"auth_type": auth.AuthType,
		})
	if result.Error != nil {
		return result.Error
	}

	// // Invalidate cache cũ và cache mới
	go func() {
		r.invalidateAuthCache(context.Background(), oldAuth.Email, oldAuth.UserId)
		r.cacheAuth(context.Background(), auth)
	}()

	return nil
}
