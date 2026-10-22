package usecase

import (
	"context"
	"errors"
	"fmt"
	domainAuth "github.com/DVV-15324/witches-auth/internal/shared/domain"
	domainToken "github.com/DVV-15324/witches-auth/internal/shared/domain"
	utils "github.com/DVV-15324/witches-auth/internal/shared/utils"
	pkg "github.com/DVV-15324/witches-auth/pkg"
	pkg_dpop "github.com/DVV-15324/witches-auth/pkg/dpop"
	"time"

	w_resp "github.com/DVV-15324/witches/pkg/core/response"

	"github.com/google/uuid"
)

func (a *AuthUseCase) Login(ctx context.Context, au *domainAuth.Auth, jwk *pkg_dpop.JWK, ipAddress string, userAgent string, timeZone string, locale string) (*pkg.TokenResponse, *w_resp.AppError) {
	auth, err := a.AuthReponsitory.GetAuthByEmail(ctx, au.Email)
	if err != nil {
		resp := w_resp.NewAppError(500, err, time.Now())
		return nil, resp
	}
	if auth == nil {
		resp := w_resp.NewAppError(500, errors.New("email không tồn tại"), time.Now())
		return nil, resp
	}

	if auth.Banned {
		resp := w_resp.NewAppError(500, errors.New("tài khoản đã bị khóa"), time.Now())
		return nil, resp
	}

	if !a.Core.Hash.CompareHashAndPassword(auth.Password, au.Password, auth.Salt) {

		resp := w_resp.NewAppError(500, errors.New("email hoặc mật khẩu không đúng"), time.Now())
		return nil, resp
	}

	s := utils.NewUID(auth.UserId, 1)
	sub := s.ToBase58()
	tid := uuid.New().String()
	jkt, err_jkt := pkg_dpop.JWKThumbprint(jwk)
	if err_jkt != nil {
		e := fmt.Sprintf("JWKThumbprint error: %v", err_jkt)
		resp := w_resp.NewAppError(500, errors.New(e), time.Now())
		return nil, resp
	}

	// revoke
	_, _ = a.RefreshTokenUseCase.RevokeByUserAndJKT(ctx, auth.UserId, jkt, "re-login")

	accessToken, err_access_token := a.Core.JWT.IssueAccessToken(ctx, sub, tid, jkt)
	if err_access_token != nil {
		resp := w_resp.NewAppError(500, err_access_token, time.Now())
		return nil, resp
	}
	refreshTokenStr, err_refresh_token := a.Core.JWT.IssueRefreshToken(ctx, sub, tid, jkt)
	if err_refresh_token != nil {
		resp := w_resp.NewAppError(500, err_refresh_token, time.Now())
		return nil, resp
	}

	refreshToken := &domainToken.RefreshToken{
		UserID:    auth.UserId,
		Token:     refreshTokenStr.Token,
		JKT:       jkt,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		Locale:    locale,
		Timezone:  timeZone,
		ExpiresAt: time.Now().Add(time.Duration(a.Core.JWT.RefreshTokenExpiry) * time.Second),
	}
	if err := a.RefreshTokenUseCase.Create(ctx, refreshToken); err != nil {
		return nil, err
	}

	session := &pkg.SessionCache{
		SessionID:   uuid.NewString(),
		UserID:      auth.UserId,
		JKT:         jkt,
		LastActive:  time.Now().Unix(),
		AccessToken: accessToken.Token,
	}
	_ = a.Core.Session.CreateSession(ctx, session)

	return &pkg.TokenResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshTokenStr,
	}, nil
}
