package usecase

import (
	"context"
	"errors"
	"strconv"
	"time"

	domainRefresh "github.com/DVV-15324/witches-auth/internal/shared/domain"

	w_resp "github.com/DVV-15324/witches/pkg/core/response"
)

func (uc *RefreshUseCase) Create(
	ctx context.Context,
	token *domainRefresh.RefreshToken,
) *w_resp.AppError {
	if token == nil {
		return w_resp.NewAppError(400, errors.New("token cannot be nil"), time.Now())
	}
	if token.UserID == 0 {
		return w_resp.NewAppError(400, errors.New("user_id is required"), time.Now())
	}
	if token.Token == "" {
		return w_resp.NewAppError(400, errors.New("token is required"), time.Now())
	}
	if token.JKT == "" {
		return w_resp.NewAppError(400, errors.New("jkt is required"), time.Now())
	}
	REFRESH_TOKEN_TTL, _ := strconv.ParseInt(uc.Core.Config.Others["REFRESH_TOKEN_TTL"], 10, 64)
	// Set default ExpiresAt nếu chưa có
	if token.ExpiresAt.IsZero() {
		ttl := time.Duration(REFRESH_TOKEN_TTL) * time.Second
		token.ExpiresAt = time.Now().Add(ttl)
	}

	// Force state
	token.Revoked = false
	token.RevokedAt = nil
	token.RevokedReason = ""

	// Set timestamps nếu chưa có (GORM tự set nếu dùng autoCreateTime)
	now := time.Now()
	if token.CreatedAt.IsZero() {
		token.CreatedAt = now
	}
	token.UpdatedAt = now

	if err := uc.RefreshTokenRepo.Create(ctx, token); err != nil {
		return w_resp.NewAppError(500, err, time.Now())
	}

	return nil
}
