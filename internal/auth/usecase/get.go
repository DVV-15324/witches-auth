package usecase

import (
	"context"

	domainAuth "github.com/DVV-15324/witches-auth/internal/shared/domain"
	"time"

	w_resp "github.com/DVV-15324/witches/pkg/core/response"
)

func (a *AuthUseCase) GetAuthByEmail(ctx context.Context, email string) (*domainAuth.Auth, *w_resp.AppError) {

	users, err := a.AuthReponsitory.GetAuthByEmail(ctx, email)
	if err != nil {
		resp := w_resp.NewAppError(500, err, time.Now())
		return nil, resp
	}
	return users, nil
}

func (a *AuthUseCase) GetAuthByUserId(ctx context.Context, userId int64) (*domainAuth.Auth, *w_resp.AppError) {
	users, err := a.AuthReponsitory.GetAuthByUserId(ctx, userId)
	if err != nil {
		resp := w_resp.NewAppError(500, err, time.Now())
		return nil, resp
	}
	return users, nil
}
