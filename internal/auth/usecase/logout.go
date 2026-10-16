package usecase

import (
	"context"
	"errors"
	"time"

	w_resp "github.com/DVV-15324/witches/pkg/core/response"
)

func (a *AuthUseCase) Logout(ctx context.Context, accessToken, refreshTokenStr string) *w_resp.AppError {
	token, err := a.RefreshTokenUseCase.GetByToken(ctx, refreshTokenStr)
	if err != nil {
		return err
	}
	if token == nil {
		return w_resp.NewAppError(500, errors.New("refresh token không hợp lệ"), time.Now())
	}
	if err := a.RefreshTokenUseCase.Revoke(ctx, refreshTokenStr, "user logout", token.Token); err != nil {
		return err
	}
	a.Core.Blacklist.BlacklistToken(ctx, accessToken)
	a.Core.Session.DeleteSession(ctx, token.UserID, token.Token)

	return nil
}
