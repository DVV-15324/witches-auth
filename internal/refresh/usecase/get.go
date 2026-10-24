package usecase

import (
	"context"
	"errors"
	domainRefresh "github.com/DVV-15324/witches-auth/pkg/shared/domain"
	"time"

	w_resp "github.com/DVV-15324/witches/pkg/core/response"
)

// GetByToken - Lấy refresh token theo token string
func (uc *RefreshUseCase) GetByToken(ctx context.Context, token string) (*domainRefresh.RefreshToken, *w_resp.AppError) {
	if token == "" {
		return nil, w_resp.NewAppError(400, errors.New("token cannot be empty"), time.Now())
	}
	// Get from repository (có cache)
	storedToken, err := uc.RefreshTokenRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, w_resp.NewAppError(500, err, time.Now())
	}
	if storedToken == nil {
		return nil, w_resp.NewAppError(404, errors.New("token not found"), time.Now())
	}

	return storedToken, nil
}

// GetByJKT - Lấy refresh token theo JKT string
func (uc *RefreshUseCase) GetByJKT(ctx context.Context, JKT string) (*domainRefresh.RefreshToken, *w_resp.AppError) {
	if JKT == "" {
		return nil, w_resp.NewAppError(400, errors.New("JKT cannot be empty"), time.Now())
	}
	// Get from repository (có cache)
	storedToken, err := uc.RefreshTokenRepo.GetByJKT(ctx, JKT)
	if err != nil {
		return nil, w_resp.NewAppError(500, err, time.Now())
	}
	if storedToken == nil {
		return nil, w_resp.NewAppError(404, errors.New("JKT not found"), time.Now())
	}

	return storedToken, nil
}
