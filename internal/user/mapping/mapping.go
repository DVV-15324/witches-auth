package mapping

import (
	"github.com/DVV-15324/witches-auth/internal/shared/utils"
	dtoUser "github.com/DVV-15324/witches-auth/internal/user/dto/response"
	modelUser "github.com/DVV-15324/witches-auth/internal/user/model"
	domainUser "github.com/DVV-15324/witches-auth/pkg/shared/domain"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
)

func FromDtoToDomainUser(dto *dtoUser.User) *domainUser.User {
	if dto == nil {
		return nil
	}
	return &domainUser.User{
		ID:        utils.DecodeFromBase58(dto.ID).LocalID,
		Name:      dto.Name,
		Email:     dto.Email,
		Role:      dto.Role,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}
func FromDtoToDomainUserList(dtos []*dtoUser.User) []*domainUser.User {
	return w_utils.MapPtrSlice(dtos, FromDtoToDomainUser)
}

func FromDomainToModelUser(domain *domainUser.User) *modelUser.User {
	if domain == nil {
		return nil
	}
	return &modelUser.User{
		ID:        domain.ID,
		Name:      domain.Name,
		Email:     domain.Email,
		Role:      domain.Role,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func FromDomainToModelUserList(domains []*domainUser.User) []*modelUser.User {
	return w_utils.MapPtrSlice(domains, FromDomainToModelUser)
}

func FromModelToDomainUser(model *modelUser.User) *domainUser.User {
	if model == nil {
		return nil
	}
	return &domainUser.User{
		ID:        model.ID,
		Name:      model.Name,
		Email:     model.Email,
		Role:      model.Role,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func FromModelToDomainUserList(entities []*modelUser.User) []*domainUser.User {
	return w_utils.MapPtrSlice(entities, FromModelToDomainUser)
}

func FromDomainToDtoUser(domain *domainUser.User) *dtoUser.User {
	if domain == nil {
		return nil
	}
	uid := utils.NewUID((domain.ID), utils.ObjectUser)
	if err := uid.Validate(); err != nil {
		return nil
	}

	return &dtoUser.User{
		ID:        uid.ToBase58(),
		Name:      domain.Name,
		Email:     domain.Email,
		Role:      domain.Role,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func FromDomainToDtoUserList(domains []*domainUser.User) []*dtoUser.User {
	return w_utils.MapPtrSlice(domains, FromDomainToDtoUser)
}
