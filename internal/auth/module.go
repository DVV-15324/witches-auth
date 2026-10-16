package auth

import (
	core "example/cmd/server/core"
	dtoAuth "example/internal/auth/dto/request"
	handler_auth "example/internal/auth/handler"
	"example/internal/auth/repository"
	usecase_auth "example/internal/auth/usecase"
	usecase_refresh "example/internal/refresh/usecase"
	usecase_user "example/internal/user/usecase"
	w_handl "github.com/DVV-15324/witches/pkg/core/handler"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	"github.com/ulule/limiter/v3"
)

type AuthModule struct {
	Core    *core.CoreServices
	Handler *handler_auth.AuthHandle
	Usecase *usecase_auth.AuthUseCase
	Repo    *repository.AuthRepository
}

func NewAuthModule(
	core *core.CoreServices,
	userUsecase *usecase_user.UserUseCase,
	refreshUsecase *usecase_refresh.RefreshUseCase,
) *AuthModule {
	repo := repository.NewAuthRepository(core)
	usecase := usecase_auth.NewAuthUseCase(core, repo, userUsecase, refreshUsecase)
	handler := handler_auth.NewAuthHandle(core, usecase)
	return &AuthModule{
		Core:    core,
		Handler: handler,
		Usecase: usecase,
		Repo:    repo,
	}
}

func (m *AuthModule) RegisterPublicRoutes(gen *w_handl.SwaggerGenerator) {
	gen.POST("/v1/auth/login").
		Summary("User login").
		Description("Authenticate user and return JWT token").
		Tags("auth").
		HeaderParam("User-Agent", "User-Agent", false).
		HeaderParam("Accept-Language", "Preferred language (e.g., vi-VN, en-US)", false).
		HeaderParam("X-Timezone", "Timezone (e.g., Asia/Ho_Chi_Minh)", false).
		RateLimit(limiter.Rate{Period: m.Core.Limiter.Rate.Period, Limit: m.Core.Limiter.Rate.Limit}).
		Body(dtoAuth.Login{}, "Login credentials").
		Response(200, w_resp.AppResponse{}, "Login successful").
		Response(401, w_resp.AppResponse{}, "Unauthorized").
		Handler(m.Handler.HandleLogin()).
		Build()

	gen.POST("/v1/auth/register").
		Summary("User registration").
		Description("Register a new user account").
		Tags("auth").
		HeaderParam("User-Agent", "User-Agent", false).
		HeaderParam("Accept-Language", "Preferred language (e.g., vi-VN, en-US)", false).
		HeaderParam("X-Timezone", "Timezone (e.g., Asia/Ho_Chi_Minh)", false).
		RateLimit(limiter.Rate{Period: m.Core.Limiter.Rate.Period, Limit: m.Core.Limiter.Rate.Limit}).
		Body(dtoAuth.Register{}, "User registration data").
		Response(201, w_resp.AppResponse{}, "Registration successful").
		Response(400, w_resp.AppResponse{}, "Bad request").
		Handler(m.Handler.HandleRegister()).
		Build()

}

func (m *AuthModule) RegisterProtectedRoutes(gen *w_handl.SwaggerGenerator) {
	gen.POST("/v1/auth/logout").
		Summary("User logout").
		Description("Logout user and invalidate JWT token").
		Tags("auth").
		HeaderParam("User-Agent", "User-Agent", false).
		HeaderParam("Accept-Language", "Preferred language (e.g., vi-VN, en-US)", false).
		HeaderParam("X-Timezone", "Timezone (e.g., Asia/Ho_Chi_Minh)", false).
		RateLimit(limiter.Rate{Period: m.Core.Limiter.Rate.Period, Limit: m.Core.Limiter.Rate.Limit}).
		HeaderParam("Authorization", "DPoP access token for authentication", false).
		HeaderParam("DPoP", "DPoP proof JWT", false).
		Response(200, w_resp.AppResponse{}, "Logout successful").
		Response(401, w_resp.AppResponse{}, "Unauthorized").
		Use(m.Core.Middleware).
		Handler(m.Handler.HandleLogout()).
		Build()
}
