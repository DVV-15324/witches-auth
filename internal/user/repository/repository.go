package repository

import (
	"context"
	"example/cmd/server/core"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	"gorm.io/gorm"
)

type UserRepository struct {
	core *core.CoreServices
}

func NewUserRepository(core *core.CoreServices) *UserRepository {
	return &UserRepository{
		core: core,
	}
}

func (u *UserRepository) getDB(ctx context.Context) *gorm.DB {
	tx, err := w_utils.GetTxFromContext(ctx)
	if err == nil {
		return tx
	}
	return u.core.DB
}

