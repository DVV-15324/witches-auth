package repository

import (
	"context"
	domainUser "example/internal/shared/domain"
	mapping "example/internal/user/mapping"
)

func (u *UserRepository) CreateUser(ctx context.Context, user *domainUser.User) (int64, error) {
	modelU := mapping.FromDomainToModelUser(user)
	result := u.core.DB.WithContext(ctx).Select("email", "name").Create(modelU)
	if result.Error != nil {
		return 0, result.Error
	}
	user.ID = modelU.ID
	return user.ID, nil
}
