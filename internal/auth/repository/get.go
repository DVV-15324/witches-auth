package repository

import (
	"context"
	"errors"
	"example/internal/auth/mapping"
	modelAuth "example/internal/auth/model"
	domainAuth "example/internal/shared/domain"

	"gorm.io/gorm"
)

func (r *AuthRepository) GetAuthByEmail(ctx context.Context, email string) (*domainAuth.Auth, error) {
	auth, err := r.getAuthFromCacheByEmail(ctx, email)
	if err == nil && auth != nil {
		domainA := mapping.FromModelToDomainAuth(auth)
		return domainA, nil
	}
	var data modelAuth.Auth
	err = r.core.DB.WithContext(ctx).Where("email = ?", email).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err

	}
	r.cacheAuth(ctx, &data)
	domainA := mapping.FromModelToDomainAuth(&data)
	return domainA, nil
}
func (r *AuthRepository) GetAuthByUserId(ctx context.Context, userID int64) (*domainAuth.Auth, error) {
	auth, err := r.getAuthFromCacheByUserID(ctx, userID)
	if err == nil && auth != nil {
		domainA := mapping.FromModelToDomainAuth(auth)
		return domainA, nil
	} else {
		var data modelAuth.Auth
		err = r.core.DB.WithContext(ctx).Where("user_id = ?", userID).First(&data).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {

				return nil, nil
			}
			return nil, err
		}
		r.cacheAuth(ctx, &data)
		domainA := mapping.FromModelToDomainAuth(&data)
		return domainA, nil
	}

}
