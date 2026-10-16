package repository

import (
	"example/cmd/server/core"
)

type RefreshTokenRepository struct {
	core *core.CoreServices
}

func NewRefreshTokenRepository(core *core.CoreServices) *RefreshTokenRepository {
	return &RefreshTokenRepository{core: core}
}
