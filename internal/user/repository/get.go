package repository

import (
	"context"
	"errors"
	domainUser "example/internal/shared/domain"
	mapping "example/internal/user/mapping"
	modelUser "example/internal/user/model"

	w_utils "github.com/DVV-15324/witches/pkg/core/utils"

	"gorm.io/gorm"
)

func (u *UserRepository) GetAllUser(ctx context.Context, req *w_utils.PaginationRequest) ([]*domainUser.User, int64, error) {
	var users []*modelUser.User
	var total int64

	query := u.core.DB.WithContext(ctx).Model(&modelUser.User{})

	// Search filter
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("LOWER(name) ? OR LOWER(email) ILIKE ?", search, search)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sort := req.GetSort("created_at")
	order := req.GetOrder("desc")
	query = query.Order(sort + " " + order)

	// Apply pagination
	if err := query.Scopes(req.GormScope()).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	//mapping
	domainsMapping := mapping.FromModelToDomainUserList(users)
	return domainsMapping, total, nil
}

func (u *UserRepository) GetUserById(ctx context.Context, id int64) (*domainUser.User, error) {
	var data modelUser.User

	err := u.core.DB.WithContext(ctx).First(&data, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	domainU := mapping.FromModelToDomainUser(&data)
	return domainU, nil
}
