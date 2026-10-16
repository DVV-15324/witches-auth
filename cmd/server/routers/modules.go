package routers

import (
	"example/cmd/server/core"
	"example/internal/auth"
	"example/internal/refresh"
	"example/internal/user"
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
