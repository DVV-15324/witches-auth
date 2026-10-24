package usecase

import (
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	"time"

	"context"
	domainUser "github.com/DVV-15324/witches-auth/pkg/shared/domain"
)

func (u *UserUseCase) CreateUser(ctx context.Context, user *domainUser.User) (int64, *w_resp.AppError) {
	uid, err := u.UserResponsitory.CreateUser(ctx, user)
	if err != nil {
		resp := w_resp.NewAppError(500, err, time.Now())
		return 0, resp
	}
	return uid, nil
}
