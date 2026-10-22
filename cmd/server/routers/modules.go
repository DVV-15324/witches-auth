package routers

import (
	"github.com/DVV-15324/witches-auth/cmd/server/core"
	"github.com/DVV-15324/witches-auth/internal/auth"
	"github.com/DVV-15324/witches-auth/internal/refresh"
	"github.com/DVV-15324/witches-auth/internal/user"
)

type Modules struct {
	Auth    *auth.AuthModule
	User    *user.UserModule
	Refresh *refresh.RefreshModule
}

func InitModules(core *core.CoreServices) *Modules {
	userModule := user.NewUserModule(core)
	refreshModule := refresh.NewRefreshModule(core, userModule.Usecase)
	authModule := auth.NewAuthModule(core, userModule.Usecase, refreshModule.Usecase)
	refreshModule.SetAuthUseCase(authModule.Usecase)

	return &Modules{
		Auth:    authModule,
		User:    userModule,
		Refresh: refreshModule,
	}
}
