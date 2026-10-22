package usecase

import (
	"context"
	"github.com/DVV-15324/witches-auth/cmd/server/core"
	domainAuth "github.com/DVV-15324/witches-auth/internal/shared/domain"
	domainRefresh "github.com/DVV-15324/witches-auth/internal/shared/domain"
	domainUser "github.com/DVV-15324/witches-auth/internal/shared/domain"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
)

type IRefreshTokenRepository interface {
	Create(ctx context.Context, token *domainRefresh.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*domainRefresh.RefreshToken, error)
	Revoke(ctx context.Context, token, reason string) error
	//RevokeAllByUser(ctx context.Context, userID int64, reason string) ([]string, error)
	RevokeByUserAndJKT(ctx context.Context, userID int64, jkt, reason string) ([]string, error)
	GetByJKT(ctx context.Context, JKT string) (*domainRefresh.RefreshToken, error)
}

type IAuthUsecase interface {
	GetAuthByUserId(ctx context.Context, userId int64) (*domainAuth.Auth, *w_resp.AppError)
	GetAuthByEmail(ctx context.Context, email string) (*domainAuth.Auth, *w_resp.AppError)
}

type IUserUseCase interface {
	CreateUser(ctx context.Context, user *domainUser.User) (int64, *w_resp.AppError)
	GetUserById(ctx context.Context, id int64) (*domainUser.User, *w_resp.AppError)
}

type RefreshUseCase struct {
	Core             *core.CoreServices
	AuthUsecase      IAuthUsecase
	RefreshTokenRepo IRefreshTokenRepository
	UserUseCase      IUserUseCase
}

func NewRefreshUseCase(core *core.CoreServices, userUseCase IUserUseCase, refreshTokenRepo IRefreshTokenRepository) *RefreshUseCase {
	return &RefreshUseCase{
		UserUseCase:      userUseCase,
		RefreshTokenRepo: refreshTokenRepo,
		Core:             core,
	}
}

func (r *RefreshUseCase) SetAuthUseCase(authUsecase IAuthUsecase) {
	r.AuthUsecase = authUsecase
}
