package usecase

import (
	"context"
	core "github.com/DVV-15324/witches-auth/cmd/server/core"
	domainAuth "github.com/DVV-15324/witches-auth/internal/shared/domain"
	domainRefresh "github.com/DVV-15324/witches-auth/internal/shared/domain"
	domainUser "github.com/DVV-15324/witches-auth/internal/shared/domain"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
)

type IAuthReponsitory interface {
	CreateAuth(ctx context.Context, auth *domainAuth.Auth) error
	GetAuthByEmail(ctx context.Context, email string) (*domainAuth.Auth, error)
	GetAuthByUserId(ctx context.Context, userID int64) (*domainAuth.Auth, error)
}

type IRefreshTokenUsecase interface {
	Create(ctx context.Context, token *domainRefresh.RefreshToken) *w_resp.AppError
	GetByToken(ctx context.Context, token string) (*domainRefresh.RefreshToken, *w_resp.AppError)
	Revoke(ctx context.Context, refreshToken string, reason string, jkt string) *w_resp.AppError
	RevokeByUserAndJKT(ctx context.Context, userID int64, jkt string, reason string) ([]string, *w_resp.AppError)
}

type IUserUseCase interface {
	CreateUser(ctx context.Context, user *domainUser.User) (int64, *w_resp.AppError)
	GetUserById(ctx context.Context, id int64) (*domainUser.User, *w_resp.AppError)
}

type AuthUseCase struct {
	Core                *core.CoreServices
	UserUseCase         IUserUseCase
	AuthReponsitory     IAuthReponsitory
	RefreshTokenUseCase IRefreshTokenUsecase
}

func NewAuthUseCase(
	core *core.CoreServices,
	authReponsitory IAuthReponsitory,
	userUseCase IUserUseCase,
	refreshTokenUseCase IRefreshTokenUsecase,
) *AuthUseCase {
	return &AuthUseCase{
		Core:                core,
		AuthReponsitory:     authReponsitory,
		UserUseCase:         userUseCase,
		RefreshTokenUseCase: refreshTokenUseCase,
	}
}
