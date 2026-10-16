package handler

import (
	"context"
	core "example/cmd/server/core"
	pkg "example/pkg"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
)

type IRefreshUseCase interface {
	Refresh(ctx context.Context, refreshTokenStr, deviceID string) (*pkg.Token, *w_resp.AppError)
}
type RefreshHandle struct {
	RefreshUsecase IRefreshUseCase
	Core           *core.CoreServices
}

func NewRefreshHandle(core *core.CoreServices, refreshseCase IRefreshUseCase) *RefreshHandle {
	return &RefreshHandle{
		RefreshUsecase: refreshseCase,
		Core:           core,
	}
}
