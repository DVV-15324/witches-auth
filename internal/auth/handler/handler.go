package handler

import (
	"context"
	"example/cmd/server/core"
	domainAuth "example/internal/shared/domain"
	pkg "example/pkg"
	pkg_dpop "example/pkg/dpop"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
)

type IUsecaseAuth interface {
	Login(ctx context.Context, au *domainAuth.Auth, jwk *pkg_dpop.JWK, ipAddress string, userAgent string, timeZone string, locale string) (*pkg.TokenResponse, *w_resp.AppError)
	Register(ctx context.Context, auth *domainAuth.Auth, name string) *w_resp.AppError
	Logout(ctx context.Context, accessToken, refreshTokenStr string) *w_resp.AppError
}
type AuthHandle struct {
	Core        *core.CoreServices
	UsecaseAuth IUsecaseAuth
}

func NewAuthHandle(core *core.CoreServices, authUseCase IUsecaseAuth) *AuthHandle {
	return &AuthHandle{
		UsecaseAuth: authUseCase,
		Core:        core,
	}
}
