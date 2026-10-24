package handler

import (
	"context"
	"github.com/DVV-15324/witches-auth/cmd/server/core"
	pkg "github.com/DVV-15324/witches-auth/pkg"
	pkg_dpop "github.com/DVV-15324/witches-auth/pkg/dpop"
	domainAuth "github.com/DVV-15324/witches-auth/pkg/shared/domain"
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
