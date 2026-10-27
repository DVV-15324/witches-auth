package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	w_resp "github.com/DVV-15324/witches/pkg/core/response"
)

// Sentinel errors — test bằng errors.Is
var (
	ErrTokenEmpty     = errors.New("token cannot be empty")
	ErrJKTRequired    = errors.New("jkt required")
	ErrUserIDRequired = errors.New("user_id required")
	ErrForbidden      = errors.New("forbidden")
	ErrTokenNotFound  = errors.New("token not found")
)

// Revoke — thu hồi 1 refresh token (logout 1 device).
func (uc *RefreshUseCase) Revoke(
	ctx context.Context,
	refreshToken string,
	jkt string,
	reason string,
) *w_resp.AppError {
	if refreshToken == "" {
		return w_resp.NewAppError(400, ErrTokenEmpty, time.Now())
	}
	if jkt == "" {
		return w_resp.NewAppError(400, ErrJKTRequired, time.Now())
	}

	storedToken, err := uc.RefreshTokenRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		return w_resp.NewAppError(500, err, time.Now())
	}
	if storedToken == nil {
		return w_resp.NewAppError(404, ErrTokenNotFound, time.Now())
	}

	// Verify jkt trước — chống dùng token của device khác
	if storedToken.JKT != jkt {
		return w_resp.NewAppError(403, ErrForbidden, time.Now())
	}

	// Idempotent: nếu chưa revoke thì revoke; nếu rồi thì bỏ qua bước DB
	if !storedToken.Revoked {
		if err := uc.RefreshTokenRepo.Revoke(ctx, refreshToken, reason); err != nil {
			return w_resp.NewAppError(500, err, time.Now())
		}
	}

	// // Blacklist — quan trọng, không nuốt error
	// if err := uc.Core.Blacklist.BlacklistToken(ctx, refreshToken); err != nil {
	// 	return w_resp.NewAppError(500, fmt.Errorf("blacklist token: %w", err), time.Now())
	// }

	// Xóa session Redis
	if err := uc.Core.Session.DeleteSession(ctx, storedToken.UserID, jkt); err != nil {
		// Session có thể không tồn tại → không fail cả logout
		// Nếu cần strict hơn, đổi thành return AppError
		_ = err
	}

	return nil
}

// RevokeByUserAndJKT — thu hồi tất cả token của (userID, jkt).
// Dùng trong login: giải phóng slot trước khi insert token mới.
// Trả về danh sách token đã revoke để caller (nếu cần) blacklist thêm.
func (uc *RefreshUseCase) RevokeByUserAndJKT(
	ctx context.Context,
	userID int64,
	jkt string,
	reason string,
) ([]string, *w_resp.AppError) {
	if userID == 0 {
		return nil, w_resp.NewAppError(400, ErrUserIDRequired, time.Now())
	}
	if jkt == "" {
		return nil, w_resp.NewAppError(400, ErrJKTRequired, time.Now())
	}

	// Repo làm hết: transaction + lock + update + clear cache
	revokedTokens, err := uc.RefreshTokenRepo.RevokeByUserAndJKT(ctx, userID, jkt, reason)
	if err != nil {
		return nil, w_resp.NewAppError(500,
			fmt.Errorf("revoke by user+jkt: %w", err), time.Now())
	}

	// // Blacklist từng token vừa revoke
	// for _, token := range revokedTokens {
	// 	if err := uc.Core.Blacklist.BlacklistToken(ctx, token); err != nil {
	// 		return revokedTokens, w_resp.NewAppError(500,
	// 			fmt.Errorf("blacklist token: %w", err), time.Now())
	// 	}
	// }

	// Xóa session Redis — best effort
	if err := uc.Core.Session.DeleteSession(ctx, userID, jkt); err != nil {
		_ = err
	}

	return revokedTokens, nil
}
