package usecase

import (
	"context"
	domainUser "github.com/DVV-15324/witches-auth/internal/shared/domain"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	"time"
)

func (u *UserUseCase) GetAllUser(ctx context.Context, req *w_utils.PaginationRequest) ([]*domainUser.User, int64, *w_resp.AppError) {
	users, total, err := u.UserResponsitory.GetAllUser(ctx, req)
	if err != nil {
		resp := w_resp.NewAppError(500, err, time.Now())
		return nil, 0, resp
	}
	return users, total, nil
}

func (u *UserUseCase) GetUserById(ctx context.Context, id int64) (*domainUser.User, *w_resp.AppError) {
	user, err := u.UserResponsitory.GetUserById(ctx, id)
	if err != nil {
		resp := w_resp.NewAppError(404, err, time.Now())
		return nil, resp
	}
	return user, nil
}
