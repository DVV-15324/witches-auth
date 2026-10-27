package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/DVV-15324/witches-auth/internal/shared/utils"
	pkg "github.com/DVV-15324/witches-auth/pkg"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	"github.com/google/uuid"
)

func (uc *RefreshUseCase) Refresh(ctx context.Context, refreshTokenStr, jkt string) (*pkg.Token, *w_resp.AppError) {
	// 1. Lấy stored token
	storedToken, err := uc.RefreshTokenRepo.GetByToken(ctx, refreshTokenStr)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("GetByToken failed", "err", err)
		return nil, w_resp.NewAppError(500, errors.New("internal server error"), time.Now())
	}
	if storedToken == nil {
		return nil, w_resp.NewAppError(401, errors.New("refresh token không hợp lệ"), time.Now())
	}

	// 2. Revoked
	if storedToken.Revoked {
		return nil, w_resp.NewAppError(401, errors.New("refresh token đã bị thu hồi"), time.Now())
	}

	// 3. Hết hạn
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, w_resp.NewAppError(401, errors.New("refresh token đã hết hạn"), time.Now())
	}

	// 4. Device binding
	if storedToken.JKT != jkt {
		uc.RefreshTokenRepo.Revoke(ctx, refreshTokenStr, "Jkt mismatch")
		return nil, w_resp.NewAppError(401, errors.New("Jkt không khớp"), time.Now())
	}

	// 5. Auth
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

	// 6. Issue AT
	s := utils.NewUID(auth.UserId, 1)
	sub := s.ToBase58()
	tid := uuid.New().String()

	accessToken, err := uc.Core.JWT.IssueAccessToken(ctx, sub, tid, jkt, "")
	if err != nil {
		slog.Error("IssueAccessToken failed", "err", err, "userID", auth.UserId)
		return nil, w_resp.NewAppError(500, errors.New("internal server error"), time.Now())
	}

	// 7. Session
	if jkt != "" {
		// Mobile: key cố định theo jkt
		// if err := uc.Core.Session.UpdateSessionB(ctx, auth.UserId, "", accessToken.Token); err != nil {
		// 	slog.Warn("UpdateSession (mobile) failed", "err", err, "userID", auth.UserId)
		// }
	} else {
		// Web: tạo session mới theo sha256(newAT); old tự hết hạn qua TTL
		if err := uc.Core.Session.CreateSession(ctx, &pkg.SessionCache{
			UserID:      auth.UserId,
			JKT:         "",
			AccessToken: accessToken.Token,
			LastActive:  time.Now().Unix(),
		}); err != nil {
			slog.Warn("CreateSession (web) failed", "err", err, "userID", auth.UserId)
		}
	}

	return &pkg.Token{Token: accessToken.Token}, nil
}
