package repository

import (
	"context"
	"example/cmd/server/core"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	"gorm.io/gorm"
)

type AuthRepository struct {
	core *core.CoreServices
}

func NewAuthRepository(core *core.CoreServices) *AuthRepository {
	return &AuthRepository{
		core: core,
	}
}

func (a *AuthRepository) getDB(ctx context.Context) *gorm.DB {
	tx, err := w_utils.GetTxFromContext(ctx)
	if err == nil {
		return tx // Dùng transaction
	}
	return a.core.DB // Dùng db thường
}
