package usecase

import (
	"context"
	"errors"
	"example/internal/shared/utils"
	pkg "example/pkg"
	"fmt"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	"github.com/google/uuid"
	"time"
)

func (uc *RefreshUseCase) Refresh(ctx context.Context, refreshTokenStr, jkt string) (*pkg.Token, *w_resp.AppError) {
	// 1. Lấy stored token
	storedToken, err := uc.RefreshTokenRepo.GetByToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, w_resp.NewAppError(500, err, time.Now())
	}
	if storedToken == nil {
		return nil, w_resp.NewAppError(404, errors.New("refresh token không hợp lệ"), time.Now())
	}

	// 2. Kiểm tra revoked
	if storedToken.Revoked {
		return nil, w_resp.NewAppError(401, errors.New("refresh token đã bị thu hồi"), time.Now())
	}

	// 3. Kiểm tra hết hạn (đã sửa đúng)
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, w_resp.NewAppError(401, errors.New("refresh token đã hết hạn"), time.Now())
	}

	// 4. Kiểm tra device
	if storedToken.JKT != jkt && jkt != "" {
		uc.RefreshTokenRepo.Revoke(ctx, refreshTokenStr, "Jkt mismatch")
		return nil, w_resp.NewAppError(401, errors.New("Jkt không khớp"), time.Now())
	}

	// 5. Lấy thông tin auth (đã inject đúng AuthUsecase)
	auth, errResp := uc.AuthUsecase.GetAuthByUserId(ctx, storedToken.UserID)
	if errResp != nil {
		return nil, errResp
	}
	if auth == nil {
		return nil, w_resp.NewAppError(404, errors.New("user not found"), time.Now())
	}
	if auth.Banned {
		return nil, w_resp.NewAppError(403, errors.New("user is banned"), time.Now())
	}

	// 6. Tạo sub và tid
	s := utils.NewUID(auth.UserId, 1)
	sub := s.ToBase58()
	tid := uuid.New().String()

	// 7. Issue access token (SỬA: dùng uc.Jwt hoặc uc.Core.JWT)
	accessToken, err := uc.Core.JWT.IssueAccessToken(ctx, sub, tid, jkt)
	if err != nil {
		return nil, w_resp.NewAppError(500, err, time.Now())
	}

	// 8. Update session
	if err := uc.Core.Session.UpdateSession(ctx, auth.UserId, jkt, accessToken.Token); err != nil {
		fmt.Printf("Failed to update session: %v\n", err)
	}

	// 9. Return
	return &pkg.Token{
		Token: accessToken.Token,
	}, nil
}
