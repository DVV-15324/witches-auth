package repository

import (
	"github.com/DVV-15324/witches-auth/cmd/server/core"
)

type RefreshTokenRepository struct {
	core *core.CoreServices
}

func NewRefreshTokenRepository(core *core.CoreServices) *RefreshTokenRepository {
	return &RefreshTokenRepository{core: core}
}
