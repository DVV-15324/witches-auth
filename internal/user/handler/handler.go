package handler

import (
	"github.com/DVV-15324/witches-auth/cmd/server/core"
	domain "github.com/DVV-15324/witches-auth/pkg/shared/domain"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"

	"context"
)

type IUserUseCase interface {
	GetAllUser(ctx context.Context, req *w_utils.PaginationRequest) ([]*domain.User, int64, *w_resp.AppError)
	GetUserById(ctx context.Context, id int64) (*domain.User, *w_resp.AppError)
}
type UserHandle struct {
	UserUseCase IUserUseCase
	Core        *core.CoreServices
}

func NewUserHandle(core *core.CoreServices, userUseCase IUserUseCase) *UserHandle {
	return &UserHandle{
		UserUseCase: userUseCase,
		Core:        core,
	}
}
