package usecase

import (
	"context"
	"example/cmd/server/core"

	domainUser "example/internal/shared/domain"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
)

type IUserResponsitory interface {
	CreateUser(cxt context.Context, user *domainUser.User) (int64, error)
	GetUserById(ctx context.Context, id int64) (*domainUser.User, error)
	GetAllUser(ctx context.Context, req *w_utils.PaginationRequest) ([]*domainUser.User, int64, error)
}

type UserUseCase struct {
	core             *core.CoreServices
	UserResponsitory IUserResponsitory
}

func NewUserUsecase(core *core.CoreServices, UserResponsitory IUserResponsitory) *UserUseCase {
	return &UserUseCase{
		core:             core,
		UserResponsitory: UserResponsitory,
	}
}
